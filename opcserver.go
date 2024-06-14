package opcda

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/huskar-t/opcda/com"
	"golang.org/x/sys/windows/registry"

	"golang.org/x/sys/windows"
)

type OPCServer struct {
	iServer       *com.IOPCServer
	iCommon       *com.IOPCCommon
	iItemProperty *com.IOPCItemProperties
	groups        *OPCGroups
	Name          string
	Node          string
	clientName    string
	location      com.CLSCTX

	container *com.IConnectionPointContainer
	point     *com.IConnectionPoint
	event     *ShutdownEventReceiver
	cookie    uint32
}

// Connect connect to OPC server
func Connect(progID, node string) (opcServer *OPCServer, err error) {
	location := com.CLSCTX_LOCAL_SERVER
	if !com.IsLocal(node) {
		location = com.CLSCTX_REMOTE_SERVER
	}
	var clsid *windows.GUID
	if location == com.CLSCTX_LOCAL_SERVER {
		id, err := windows.GUIDFromString(progID)
		if err != nil {
			return nil, err
		}
		clsid = &id
	} else {
		// try get clsid from server list
		clsid, err = getClsIDFromServerList(progID, node, location)
		if err != nil {
			// try get clsid from windows reg
			clsid, err = getClsIDFromReg(progID, node)
			if err != nil {
				return nil, err
			}
		}
	}
	iUnknownServer, err := com.MakeCOMObjectEx(node, location, clsid, &com.IID_IOPCServer)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			iUnknownServer.Release()
		}
	}()
	var iUnknownCommon *com.IUnknown
	err = iUnknownServer.QueryInterface(&com.IID_IOPCCommon, unsafe.Pointer(&iUnknownCommon))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			iUnknownCommon.Release()
		}
	}()
	var iUnknownItemProperties *com.IUnknown
	err = iUnknownServer.QueryInterface(&com.IID_IOPCItemProperties, unsafe.Pointer(&iUnknownItemProperties))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			iUnknownItemProperties.Release()
		}
	}()
	server := &com.IOPCServer{IUnknown: iUnknownServer}
	common := &com.IOPCCommon{IUnknown: iUnknownCommon}
	itemProperties := &com.IOPCItemProperties{IUnknown: iUnknownItemProperties}
	var iUnknownContainer *com.IUnknown
	err = iUnknownServer.QueryInterface(&com.IID_IConnectionPointContainer, unsafe.Pointer(&iUnknownContainer))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			iUnknownContainer.Release()
		}
	}()
	container := &com.IConnectionPointContainer{IUnknown: iUnknownContainer}
	point, err := container.FindConnectionPoint(&IID_IOPCShutdown)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			point.Release()
		}
	}()
	event := NewShutdownEventReceiver()
	cookie, err := point.Advise((*com.IUnknown)(unsafe.Pointer(event)))
	if err != nil {
		return nil, err
	}
	opcServer = &OPCServer{
		iServer:       server,
		iCommon:       common,
		iItemProperty: itemProperties,
		Name:          progID,
		Node:          node,
		location:      location,
		container:     container,
		point:         point,
		event:         event,
		cookie:        cookie,
	}
	opcServer.groups = NewOPCGroups(opcServer)
	return opcServer, nil
}

func getClsIDFromServerList(progID, node string, location com.CLSCTX) (*windows.GUID, error) {
	iCatInfo, err := com.MakeCOMObjectEx(node, location, &com.CLSID_OpcServerList, &com.IID_IOPCServerList2)
	if err != nil {
		return nil, err
	}
	defer iCatInfo.Release()
	sl := &com.IOPCServerList2{IUnknown: iCatInfo}
	clsid, err := sl.CLSIDFromProgID(progID)
	if err != nil {
		return nil, err
	}
	return clsid, nil
}

func getClsIDFromReg(progID, node string) (*windows.GUID, error) {
	var clsid windows.GUID
	var err error
	hKey, err := registry.OpenRemoteKey(node, registry.CLASSES_ROOT)
	if err != nil {
		return nil, err
	}
	defer hKey.Close()
	hProgIDKey, err := registry.OpenKey(hKey, progID, registry.READ)
	if err != nil {
		return nil, err
	}
	defer hProgIDKey.Close()
	hClsidKey, err := registry.OpenKey(hProgIDKey, "CLSID", registry.READ)
	if err != nil {
		return nil, err
	}
	defer hClsidKey.Close()
	clsidStr, _, err := hClsidKey.GetStringValue("")
	if err != nil {
		return nil, err
	}
	clsid, err = windows.GUIDFromString(clsidStr)
	return &clsid, err
}

type ServerInfo struct {
	ProgID       string
	ClsStr       string
	VerIndProgID string
	ClsID        *windows.GUID
}

// GetOPCServers get OPC servers from node
func GetOPCServers(node string) ([]*ServerInfo, error) {
	location := com.CLSCTX_LOCAL_SERVER
	if !com.IsLocal(node) {
		location = com.CLSCTX_REMOTE_SERVER
	}
	iCatInfo, err := com.MakeCOMObjectEx(node, location, &com.CLSID_OpcServerList, &com.IID_IOPCServerList2)
	if err != nil {
		return nil, err
	}
	cids := []windows.GUID{IID_CATID_OPCDAServer10, IID_CATID_OPCDAServer20}
	defer iCatInfo.Release()
	sl := &com.IOPCServerList2{IUnknown: iCatInfo}
	iEnum, err := sl.EnumClassesOfCateGories(cids, nil)
	if err != nil {
		return nil, err
	}
	defer iEnum.Release()
	var result []*ServerInfo
	for {
		var classID windows.GUID
		var actual uint32
		err = iEnum.Next(1, &classID, &actual)
		if err != nil {
			break
		}
		server, err := getServer(sl, &classID)
		if err != nil {
			return nil, err
		}
		result = append(result, server)
	}
	return result, nil
}

func getServer(sl *com.IOPCServerList2, classID *windows.GUID) (*ServerInfo, error) {
	progID, userType, VerIndProgID, err := sl.GetClassDetails(classID)
	if err != nil {
		return nil, fmt.Errorf("FAILED to get prog ID from class ID: %w", err)
	}
	defer func() {
		com.CoTaskMemFree(unsafe.Pointer(progID))
		com.CoTaskMemFree(unsafe.Pointer(userType))
		com.CoTaskMemFree(unsafe.Pointer(VerIndProgID))
	}()
	clsStr := classID.String()
	return &ServerInfo{
		ProgID:       windows.UTF16PtrToString(progID),
		ClsStr:       clsStr,
		ClsID:        classID,
		VerIndProgID: windows.UTF16PtrToString(VerIndProgID),
	}, nil
}

// GetLocaleID get locale ID
func (s *OPCServer) GetLocaleID() (uint32, error) {
	localeID, err := s.iCommon.GetLocaleID()
	return localeID, err
}

// GetStartTime Returns the time the server started running
func (s *OPCServer) GetStartTime() (time.Time, error) {
	status, err := s.iServer.GetStatus()
	if err != nil {
		return time.Time{}, err
	}
	return status.StartTime, nil
}

// GetCurrentTime Returns the current time from the server
func (s *OPCServer) GetCurrentTime() (time.Time, error) {
	status, err := s.iServer.GetStatus()
	if err != nil {
		return time.Time{}, err
	}
	return status.CurrentTime, nil
}

// GetLastUpdateTime Returns the last update time from the server
func (s *OPCServer) GetLastUpdateTime() (time.Time, error) {
	status, err := s.iServer.GetStatus()
	if err != nil {
		return time.Time{}, err
	}
	return status.LastUpdateTime, nil
}

// GetMajorVersion Returns the major part of the server version number
func (s *OPCServer) GetMajorVersion() (uint16, error) {
	status, err := s.iServer.GetStatus()
	if err != nil {
		return 0, err
	}
	return status.MajorVersion, nil
}

// GetMinorVersion Returns the minor part of the server version number
func (s *OPCServer) GetMinorVersion() (uint16, error) {
	status, err := s.iServer.GetStatus()
	if err != nil {
		return 0, err
	}
	return status.MinorVersion, nil
}

// GetBuildNumber Returns the build number of the server
func (s *OPCServer) GetBuildNumber() (uint16, error) {
	status, err := s.iServer.GetStatus()
	if err != nil {
		return 0, err
	}
	return status.BuildNumber, nil
}

// GetVendorInfo Returns the vendor information string for the server
func (s *OPCServer) GetVendorInfo() (string, error) {
	status, err := s.iServer.GetStatus()
	if err != nil {
		return "", err
	}
	return status.VendorInfo, nil
}

// GetServerState Returns the server’s state
func (s *OPCServer) GetServerState() (com.OPCServerState, error) {
	status, err := s.iServer.GetStatus()
	if err != nil {
		return 0, err
	}
	return status.ServerState, nil
}

// SetLocaleID set locale ID
func (s *OPCServer) SetLocaleID(localeID uint32) error {
	return s.iCommon.SetLocaleID(localeID)
}

// GetBandwidth Returns the bandwidth of the server
func (s *OPCServer) GetBandwidth() (uint32, error) {
	status, err := s.iServer.GetStatus()
	if err != nil {
		return 0, err
	}
	return status.BandWidth, nil
}

// GetOPCGroups get a collection of OPCGroup objects
func (s *OPCServer) GetOPCGroups() *OPCGroups {
	return s.groups
}

// GetServerName Returns the server name of the server that the client connected to via Connect().
func (s *OPCServer) GetServerName() string {
	return s.Name
}

// GetServerNode Returns the node name of the server that the client connected to via Connect().
func (s *OPCServer) GetServerNode() string {
	return s.Node
}

// GetClientName Returns the client name of the client
func (s *OPCServer) GetClientName() string {
	return s.clientName
}

// SetClientName Sets the client name of the client
func (s *OPCServer) SetClientName(clientName string) error {
	err := s.iCommon.SetClientName(clientName)
	if err != nil {
		return err
	}
	s.clientName = clientName
	return nil
}

type PropertyDescription struct {
	PropertyID   int32
	Description  string
	DataType     int16
	AccessRights int16
}

// CreateBrowser Creates an OPCBrowser object
func (s *OPCServer) CreateBrowser() (*OPCBrowser, error) {
	return NewOPCBrowser(s)
}

// GetErrorString Converts an error number to a readable string
func (s *OPCServer) GetErrorString(errorCode int32) (string, error) {
	return s.iCommon.GetErrorString(uint32(errorCode))
}

// QueryAvailableLocaleIDs Return the available LocaleIDs for this server/client session
func (s *OPCServer) QueryAvailableLocaleIDs() ([]uint32, error) {
	return s.iCommon.QueryAvailableLocaleIDs()
}

// QueryAvailableProperties Return a list of ID codes and Descriptions for the available properties for this ItemID
func (s *OPCServer) QueryAvailableProperties(itemID string) (pPropertyIDs []uint32, ppDescriptions []string, ppvtDataTypes []uint16, err error) {
	return s.iItemProperty.QueryAvailableProperties(itemID)
}

// GetItemProperties Return a list of the current data values for the passed ID codes.
func (s *OPCServer) GetItemProperties(itemID string, propertyIDs []uint32) (data []interface{}, errors []error, err error) {
	var errs []int32
	data, errs, err = s.iItemProperty.GetItemProperties(itemID, propertyIDs)
	if err != nil {
		return nil, nil, err
	}
	errors = s.errors(errs)
	return data, errors, nil
}

// LookupItemIDs Return a list of ItemIDs (if available) for each of the passed ID codes.
// have not tested because simulator return error
func (s *OPCServer) LookupItemIDs(itemID string, propertyIDs []uint32) ([]string, []error, error) {
	ItemIDs, errs, err := s.iItemProperty.LookupItemIDs(itemID, propertyIDs)
	if err != nil {
		return nil, nil, err
	}
	errors := s.errors(errs)
	return ItemIDs, errors, nil
}

func (s *OPCServer) errors(errs []int32) []error {
	errors := make([]error, len(errs))
	for i, e := range errs {
		if e < 0 {
			errStr, _ := s.GetErrorString(e)
			errors[i] = &OPCError{
				ErrorCode:    e,
				ErrorMessage: errStr,
			}
		}
	}
	return errors
}

// RegisterServerShutDown register server shut down event
func (s *OPCServer) RegisterServerShutDown(ch chan string) error {
	s.event.AddReceiver(ch)
	return nil
}

// Disconnect from OPC server
func (s *OPCServer) Disconnect() error {
	err := s.point.Unadvise(s.cookie)
	s.point.Release()
	s.container.Release()
	s.groups.Release()
	s.iItemProperty.Release()
	s.iCommon.Release()
	s.iServer.Release()
	return err
}
