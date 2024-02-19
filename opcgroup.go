package opcda

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/huskar-t/opcda/com"
)

type OPCGroup struct {
	parent             *OPCGroups
	groupStateMgt      *com.IOPCGroupStateMgt
	syncIO             *com.IOPCSyncIO
	asyncIO2           *com.IOPCAsyncIO2
	iCommon            *com.IOPCCommon
	clientGroupHandle  uint32
	serverGroupHandle  uint32
	groupName          string
	revisedUpdateRate  uint32
	items              *OPCItems
	callbackLock       sync.Mutex
	container          *com.IConnectionPointContainer
	point              *com.IConnectionPoint
	event              *DataEventReceiver
	cookie             uint32
	ctx                context.Context
	cancel             context.CancelFunc
	dataChangeList     []chan *DataChangeCallBackData
	readCompleteList   []chan *ReadCompleteCallBackData
	writeCompleteList  []chan *WriteCompleteCallBackData
	cancelCompleteList []chan *CancelCompleteCallBackData
}

func NewOPCGroup(
	opcGroups *OPCGroups,
	iUnknown *com.IUnknown,
	clientGroupHandle uint32,
	serverGroupHandle uint32,
	groupName string,
	revisedUpdateRate uint32,
) (*OPCGroup, error) {
	var iUnknownSyncIO *com.IUnknown
	err := iUnknown.QueryInterface(&com.IID_IOPCSyncIO, unsafe.Pointer(&iUnknownSyncIO))
	if err != nil {
		return nil, err
	}
	var iUnknownAsyncIO2 *com.IUnknown
	err = iUnknown.QueryInterface(&com.IID_IOPCAsyncIO2, unsafe.Pointer(&iUnknownAsyncIO2))
	if err != nil {
		iUnknownSyncIO.Release()
		return nil, err
	}
	var iUnknownItemMgt *com.IUnknown
	err = iUnknown.QueryInterface(&com.IID_IOPCItemMgt, unsafe.Pointer(&iUnknownItemMgt))
	if err != nil {
		iUnknownSyncIO.Release()
		iUnknownAsyncIO2.Release()
		return nil, err
	}

	o := &OPCGroup{
		parent:            opcGroups,
		groupStateMgt:     &com.IOPCGroupStateMgt{IUnknown: iUnknown},
		syncIO:            &com.IOPCSyncIO{IUnknown: iUnknownSyncIO},
		asyncIO2:          &com.IOPCAsyncIO2{IUnknown: iUnknownAsyncIO2},
		clientGroupHandle: clientGroupHandle,
		serverGroupHandle: serverGroupHandle,
		groupName:         groupName,
		revisedUpdateRate: revisedUpdateRate,
		iCommon:           opcGroups.iCommon,
	}
	o.items = NewOPCItems(o, &com.IOPCItemMgt{IUnknown: iUnknownItemMgt}, opcGroups.iCommon)
	return o, nil
}

// GetParent Returns reference to the parent OPCServer object
func (g *OPCGroup) GetParent() *OPCGroups {
	return g.parent
}

// GetName Returns the name of the group
func (g *OPCGroup) GetName() string {
	return g.groupName
}

// SetName set the name of the group
func (g *OPCGroup) SetName(name string) error {
	err := g.groupStateMgt.SetName(name)
	if err != nil {
		return err
	}
	g.groupName = name
	return nil
}

// GetIsActive Returns whether the group is active
func (g *OPCGroup) GetIsActive() bool {
	_, b, _, _, _, _, _, _, err := g.groupStateMgt.GetState()
	if err != nil {
		return false
	}
	return b
}

// SetIsActive set whether the group is active
func (g *OPCGroup) SetIsActive(isActive bool) error {
	v := com.BoolToBOOL(isActive)
	_, err := g.groupStateMgt.SetState(nil, &v, nil, nil, nil, nil)
	return err
}

// GetClientHandle get a Long value associated with the group
func (g *OPCGroup) GetClientHandle() uint32 {
	return g.clientGroupHandle
}

// SetClientHandle set a Long value associated with the group
func (g *OPCGroup) SetClientHandle(clientHandle uint32) error {
	_, err := g.groupStateMgt.SetState(nil, nil, nil, nil, nil, &clientHandle)
	if err != nil {
		return err
	}
	g.clientGroupHandle = clientHandle
	return nil
}

// GetServerHandle get the server assigned handle for the group
func (g *OPCGroup) GetServerHandle() uint32 {
	return g.serverGroupHandle
}

// GetLocaleID get the locale identifier for the group
func (g *OPCGroup) GetLocaleID() (uint32, error) {
	_, _, _, _, _, localeID, _, _, err := g.groupStateMgt.GetState()
	return localeID, err
}

// SetLocaleID set the locale identifier for the group
func (g *OPCGroup) SetLocaleID(id uint32) error {
	_, err := g.groupStateMgt.SetState(nil, nil, nil, nil, &id, nil)
	return err
}

// GetTimeBias This property provides the information needed to convert the time stamp on the data back to the local time of the device
func (g *OPCGroup) GetTimeBias() (int32, error) {
	_, _, _, timeBias, _, _, _, _, err := g.groupStateMgt.GetState()
	return timeBias, err
}

// SetTimeBias This property provides the information needed to convert the time stamp on the data back to the local time of the device
func (g *OPCGroup) SetTimeBias(timeBias int32) error {
	_, err := g.groupStateMgt.SetState(nil, nil, &timeBias, nil, nil, nil)
	return err
}

// GetDeadband A deadband is expressed as percent of full scale (legal values 0 to 100).
func (g *OPCGroup) GetDeadband() (float32, error) {
	_, _, _, _, deadband, _, _, _, err := g.groupStateMgt.GetState()
	return deadband, err
}

// SetDeadband A deadband is expressed as percent of full scale (legal values 0 to 100).
func (g *OPCGroup) SetDeadband(deadband float32) error {
	_, err := g.groupStateMgt.SetState(nil, nil, nil, &deadband, nil, nil)
	return err
}

// GetUpdateRate
// The fastest rate at which data change events may be fired. A slow process might
// cause data changes to fire at less than this rate, but they will never exceed this rate. Rate is in
// milliseconds. This property’s default depends on the value set in the OPCGroups Collection.
// Assigning a value to this property is a “request” for a new update rate. The server may not support
// that rate, so reading the property may result in a different rate (the server will use the closest rate it
// does support).
func (g *OPCGroup) GetUpdateRate() (uint32, error) {
	updateRate, _, _, _, _, _, _, _, err := g.groupStateMgt.GetState()
	return updateRate, err
}

// SetUpdateRate set the update rate
func (g *OPCGroup) SetUpdateRate(updateRate uint32) error {
	_, err := g.groupStateMgt.SetState(&updateRate, nil, nil, nil, nil, nil)
	return err
}

// OPCItems A collection of OPCItem objects
func (g *OPCGroup) OPCItems() *OPCItems {
	return g.items
}

// SyncRead This function reads the value, quality and timestamp information for one or more items in a group.
func (g *OPCGroup) SyncRead(source com.OPCDATASOURCE, serverHandles []uint32) ([]*com.ItemState, error) {
	values, errList, err := g.syncIO.Read(source, serverHandles)
	if err != nil {
		return nil, err
	}
	var errs []error
	for _, e := range errList {
		if e < 0 {
			errs = append(errs, g.getError(e))
		}
	}
	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}
	return values, nil
}

// SyncWrite Writes values to one or more items in a group
func (g *OPCGroup) SyncWrite(serverHandles []uint32, values []interface{}) ([]error, error) {
	variants := make([]com.VARIANT, len(values))
	variantWrappers := make([]*com.VariantWrapper, len(values))
	defer func() {
		for _, variant := range variantWrappers {
			err := variant.Clear()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	for i, v := range values {
		variant, err := com.NewVariant(v)
		if err != nil {
			return nil, err
		}
		variantWrappers[i] = variant
		variants[i] = *variant.Variant
	}
	errList, err := g.syncIO.Write(serverHandles, variants)
	if err != nil {
		return nil, err
	}
	errs := make([]error, len(errList))
	for i, e := range errList {
		if e < 0 {
			errs[i] = g.getError(e)
		}
	}
	return errs, nil
}

// Release Releases the resources used by the group
func (g *OPCGroup) Release() {
	if g.event != nil {
		g.point.Unadvise(g.cookie)
		g.point.Release()
		g.container.Release()
		g.event = nil
	}
	g.items.Release()
	g.groupStateMgt.Release()
	g.syncIO.Release()
	g.asyncIO2.Release()
}

type DataChangeCallBackData struct {
	TransID           uint32
	GroupHandle       uint32
	MasterQuality     int32
	MasterErr         error
	ItemClientHandles []uint32
	Values            []interface{}
	Qualities         []uint16
	TimeStamps        []time.Time
	Errors            []error
}

// RegisterDataChange Register to receive data change events
func (g *OPCGroup) RegisterDataChange(ch chan *DataChangeCallBackData) error {
	err := g.advice()
	if err != nil {
		return err
	}
	g.dataChangeList = append(g.dataChangeList, ch)
	return nil
}

// RegisterReadComplete Register to receive read complete events
func (g *OPCGroup) RegisterReadComplete(ch chan *ReadCompleteCallBackData) error {
	err := g.advice()
	if err != nil {
		return err
	}
	g.readCompleteList = append(g.readCompleteList, ch)
	return nil
}

// RegisterWriteComplete Register to receive write complete events
func (g *OPCGroup) RegisterWriteComplete(ch chan *WriteCompleteCallBackData) error {
	err := g.advice()
	if err != nil {
		return err
	}
	g.writeCompleteList = append(g.writeCompleteList, ch)
	return nil
}

// RegisterCancelComplete Register to receive cancel complete events
func (g *OPCGroup) RegisterCancelComplete(ch chan *CancelCompleteCallBackData) error {
	err := g.advice()
	if err != nil {
		return err
	}
	g.cancelCompleteList = append(g.cancelCompleteList, ch)
	return nil
}

type ReadCompleteCallBackData struct {
	TransID           uint32
	GroupHandle       uint32
	MasterQuality     int32
	MasterErr         error
	ItemClientHandles []uint32
	Values            []interface{}
	Qualities         []uint16
	TimeStamps        []time.Time
	Errors            []error
}

type WriteCompleteCallBackData struct {
	TransID           uint32
	GroupHandle       uint32
	MasterErr         error
	ItemClientHandles []uint32
	Errors            []error
}

type CancelCompleteCallBackData struct {
	TransID     uint32
	GroupHandle uint32
}

func (g *OPCGroup) advice() (err error) {
	g.callbackLock.Lock()
	defer g.callbackLock.Unlock()
	if g.event != nil {
		return nil
	}
	var iUnknownContainer *com.IUnknown
	err = g.groupStateMgt.QueryInterface(&com.IID_IConnectionPointContainer, unsafe.Pointer(&iUnknownContainer))
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			iUnknownContainer.Release()
		}
	}()
	container := &com.IConnectionPointContainer{IUnknown: iUnknownContainer}
	var point *com.IConnectionPoint
	point, err = container.FindConnectionPoint(&IID_IOPCDataCallback)
	if err != nil {
		return
	}
	dataChangeCB := make(chan *CDataChangeCallBackData, 100)
	readCB := make(chan *CReadCompleteCallBackData, 100)
	writeCB := make(chan *CWriteCompleteCallBackData, 100)
	cancelCB := make(chan *CCancelCompleteCallBackData, 100)
	event := NewDataEventReceiver(dataChangeCB, readCB, writeCB, cancelCB)
	var cookie uint32
	cookie, err = point.Advise((*com.IUnknown)(unsafe.Pointer(event)))
	if err != nil {
		return
	}
	g.ctx, g.cancel = context.WithCancel(context.Background())
	go g.loop(g.ctx, dataChangeCB, readCB, writeCB, cancelCB)
	g.container = container
	g.point = point
	g.event = event
	g.cookie = cookie
	return
}

func (g *OPCGroup) loop(ctx context.Context, dataChangeCB chan *CDataChangeCallBackData, readCB chan *CReadCompleteCallBackData, writeCB chan *CWriteCompleteCallBackData, cancelCB chan *CCancelCompleteCallBackData) {
	for {
		select {
		case <-ctx.Done():
			return
		case cbData := <-dataChangeCB:
			g.fireDataChange(cbData)
		case cbData := <-readCB:
			g.fireReadComplete(cbData)
		case cbData := <-writeCB:
			g.fireWriteComplete(cbData)
		case cbData := <-cancelCB:
			g.fireCancelComplete(cbData)
		}
	}
}

func (g *OPCGroup) fireDataChange(cbData *CDataChangeCallBackData) {
	masterError := error(nil)
	if (cbData.MasterErr) < 0 {
		masterError = g.getError(cbData.MasterErr)
	}
	itemErrors := make([]error, len(cbData.Errors))
	for i, e := range cbData.Errors {
		if e < 0 {
			itemErrors[i] = g.getError(e)
		}
	}
	data := &DataChangeCallBackData{
		TransID:           cbData.TransID,
		GroupHandle:       cbData.GroupHandle,
		MasterQuality:     cbData.MasterQuality,
		MasterErr:         masterError,
		ItemClientHandles: cbData.ItemClientHandles,
		Values:            cbData.Values,
		Qualities:         cbData.Qualities,
		TimeStamps:        cbData.TimeStamps,
		Errors:            itemErrors,
	}
	for _, backData := range g.dataChangeList {
		select {
		case backData <- data:
		default:
		}
	}
}

func (g *OPCGroup) fireReadComplete(cbData *CReadCompleteCallBackData) {
	masterError := error(nil)
	if (cbData.MasterErr) < 0 {
		masterError = g.getError(cbData.MasterErr)
	}
	itemErrors := make([]error, len(cbData.Errors))
	for i, e := range cbData.Errors {
		if e < 0 {
			itemErrors[i] = g.getError(e)
		}
	}
	data := &ReadCompleteCallBackData{
		TransID:           cbData.TransID,
		GroupHandle:       cbData.GroupHandle,
		MasterQuality:     cbData.MasterQuality,
		MasterErr:         masterError,
		ItemClientHandles: cbData.ItemClientHandles,
		Values:            cbData.Values,
		Qualities:         cbData.Qualities,
		TimeStamps:        cbData.TimeStamps,
		Errors:            itemErrors,
	}
	for _, backData := range g.readCompleteList {
		select {
		case backData <- data:
		default:
		}
	}
}

func (g *OPCGroup) fireWriteComplete(cbData *CWriteCompleteCallBackData) {
	masterError := error(nil)
	if (cbData.MasterErr) < 0 {
		masterError = g.getError(cbData.MasterErr)
	}
	itemErrors := make([]error, len(cbData.Errors))
	for i, e := range cbData.Errors {
		if e < 0 {
			itemErrors[i] = g.getError(e)
		}
	}
	data := &WriteCompleteCallBackData{
		TransID:           cbData.TransID,
		GroupHandle:       cbData.GroupHandle,
		MasterErr:         masterError,
		ItemClientHandles: cbData.ItemClientHandles,
		Errors:            itemErrors,
	}
	for _, backData := range g.writeCompleteList {
		select {
		case backData <- data:
		default:
		}
	}
}

func (g *OPCGroup) fireCancelComplete(cbData *CCancelCompleteCallBackData) {
	data := &CancelCompleteCallBackData{
		TransID:     cbData.TransID,
		GroupHandle: cbData.GroupHandle,
	}
	for _, backData := range g.cancelCompleteList {
		backData <- data
	}
}

// AsyncRead Read one or more items in a group. The results are returned via the AsyncReadComplete event associated with the OPCGroup object.
func (g *OPCGroup) AsyncRead(
	serverHandles []uint32,
	clientTransactionID uint32,
) (cancelID uint32, errs []error, err error) {
	var es []int32
	cancelID, es, err = g.asyncIO2.Read(
		serverHandles,
		clientTransactionID,
	)
	if err != nil {
		return
	}
	errs = make([]error, len(es))
	for i, e := range es {
		if e < 0 {
			errs[i] = g.getError(e)
		}
	}
	return
}

// AsyncWrite Write one or more items in a group. The results are returned via the AsyncWriteComplete event associated with the OPCGroup object.
func (g *OPCGroup) AsyncWrite(
	serverHandles []uint32,
	values []interface{},
	clientTransactionID uint32,
) (cancelID uint32, errs []error, err error) {
	variants := make([]com.VARIANT, len(values))
	variantWrappers := make([]*com.VariantWrapper, len(values))

	defer func() {
		for _, variant := range variants {
			variant.Clear()
		}
	}()
	for i, v := range values {
		variant, err := com.NewVariant(v)
		if err != nil {
			return 0, nil, err
		}
		variantWrappers[i] = variant
		variants[i] = *variant.Variant
	}
	var es []int32
	cancelID, es, err = g.asyncIO2.Write(
		serverHandles,
		variants,
		clientTransactionID,
	)
	if err != nil {
		return
	}
	errs = make([]error, len(es))
	for i, e := range es {
		if e < 0 {
			errs[i] = g.getError(e)
		}
	}
	return
}

// AsyncRefresh Generate an event for all active items in the group (whether they have changed or not). Inactive
// items are not included in the callback. The results are returned via the DataChange event
// associated with the OPCGroup object.
func (g *OPCGroup) AsyncRefresh(
	source com.OPCDATASOURCE,
	clientTransactionID uint32,
) (cancelID uint32, err error) {
	cancelID, err = g.asyncIO2.Refresh2(
		source,
		clientTransactionID,
	)
	return
}

// AsyncCancel Request that the server cancel an outstanding transaction. An AsyncCancelComplete event will
// occur indicating whether or not the cancel succeeded.
func (g *OPCGroup) AsyncCancel(cancelID uint32) error {
	return g.asyncIO2.Cancel2(cancelID)
}

func (g *OPCGroup) getError(errorCode int32) error {
	errStr, _ := g.iCommon.GetErrorString(uint32(errorCode))
	return &OPCError{
		ErrorCode:    errorCode,
		ErrorMessage: errStr,
	}
}
