package models

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	CodeOK = 0

	// 仅当返回400错误的时候才返回如下错误code
	CodeUserExists          = 30001 // 该用户已存在
	CodePasswordError       = 30002 // 账号或密码错误
	CodePasswordInvalid     = 30003 // 密码无效
	CodeUserNameInvalid     = 30004 // 用户名无效
	CodeEmailInvalid        = 30005 // email无效
	CodePhoneInvalid        = 30006 // phone无效
	CodePermissionsInvalid  = 30007 // permissions无效
	CodeProfileInvalid      = 30008 // profile无效
	CodePhoneExists         = 30011 // 用户电话号码已注册
	CodeThirdAuthExists     = 30021 // 第三方授权已存在
	CodeNotCSVFile          = 30022 // csv文件格式错误
	CodeFileContentNull     = 30023 // 文件内容为空
	CodeFileNameInvalid     = 30024 // 文件名称不合法
	CodeUserRoleForbidLogin = 30025 // 该用户无权限登录平台

	CodeNotFound         = 10001 // 未找到相关数据，请检查参数
	CodeParameterInvalid = 10002 // 参数错误
	CodeNotSupported     = 10003 // 不支持
	CodeAlreadyRecycle   = 10006 // 已回收
	CodeQueueIsNotExist  = 10008 // 白名单不存在

	CodeFeatureLimit          = 10100 // feature模板数超限制
	CodeFeatureConditionLimit = 10101 // feature条件数超限制
	CodeFeatureValuesLimit    = 10102 // feature条件填值数超限制

	CodeAlreadyUsed   = 10010 // 已被使用
	CodeAlreadyExists = 10015 // 已存在
	CodeDoesNotExist  = 10016 // 不存在

	CodeDataInvalid = 10020 // 数据无效
	CodeEmptyData   = 10021 // 空数据

	CodeSMSMany    = 20001 // 短信发送频繁请求
	CodeSMSInvalid = 20002 // 短信验证码校验失败

	CodeTriggerRuleNumOverLimit          = 30101 // 触发器规则数量超过限制
	CodeTriggerRuleConditionNumOverLimit = 30102 // 触发器规则条件数超过限制
	CodeTriggerRuleActionNumOverLimit    = 30103 // 触发器规则操作数超过限制

	CodeEventAlgorithmNumOverLimit       = 40001 // 事件算法数量超过数量限制
	CodeCountAlgorithmNumOverLimit       = 40002 // 统计算法数量超过数量限制
	CodeODBAAlgorithmNumOverLimit        = 40003 // ODBA算法数量超过数量限制
	CodeEventAlgorithmEnableNumOverLimit = 40004 // 事件算法开启数量超过限制
	CodeCountAlgorithmEnableNumOverLimit = 40005 // 统计算法开启数量超过限制
	CodeODBAAlgorithmEnableNumOverLimit  = 40006 // ODBA算法开启数量超过限制
	CodeEventAlgorithmEnableExists       = 40007 // 已启用相同类型的事件算法
	CodeUserAlgorithmEventLibOverLimit   = 40008 // 用户事件库数量超过限制

	CodeWebShareLost         = 10103 //分享已失效
	CodeWebShareExpired      = 10104 //分享已过期
	CodeWebShareBiologicAuth = 10105 //生物分享失效

	// 订单相关
	CodeOrderProductNotMatch         = 10011 // 订单所选产品有误，或是原产品发生了变更
	CodeOrderProductPlatformNotMatch = 10012 // 订单所选产品不属于该平台
	CodeOrderProductDuplicate        = 10013 // 订单所选产品重复了
	CodeOrderAddDeviceNotFound       = 10014 // 订单添加的设备，找不到
	CodeOrderAddDeviceHasAssigned    = 10015 // 订单所选产品已被分配给组织了
	CodeOrderDuplicatePay            = 10016 // 订单重复支付
	CodeOrderCompanyHasAgency        = 10017 // 创建订单所选的公司属于某代理商
	CodeOrderProductNotExists        = 10018 // 订单中有产品型号不存在
)

const (
	CodeCommon              = 50000 // 通用错误
	CodeFuncFeature         = 50001
	CodeFuncCollaboration   = 50002
	CodeProduct             = 50003 // 产品型号
	CodeCommand             = 50004 // command
	CodeSetting             = 50005 // setting
	CodeFactorySetting      = 50006 // factory setting
	CodeFile                = 50007 // file
	CodeDevice              = 50008 // device
	CodeGatewaySetting      = 50010 // 网关配置
	CodeTriggerSetting      = 50011 // 边缘智能
	CodeUser                = 50012 // 用户
	CodeFirmware            = 50014 // 固件管理
	CodeLocation            = 50016 // 定位数据
	CodeFeatureSettingFile  = 50017 // 功能配置文件
	CodeEventAlgorithm      = 50018 // 事件算法
	CodeBiological          = 50019 // 生物
	CodeVHFSetting          = 50020 // VHF配置
	CodeGeo                 = 50021 // Geo工具
	CodeCaptureImageSetting = 50022 // 图片配置
	CodeExport              = 50024 // 数据导出
	CodeRegisterTask        = 50025 // 注册任务
	CodeThird               = 50026 // third device
	CodeEntrust             = 50050 // 委托
)

const (
	//CodeCommon
	CodeCommonSystemError            = 100001 //系统错误
	CodeCommonParameterError         = 100002 //参数错误
	CodeCommonFormatTimeError        = 100003 // 时间格式化错误
	CodeCommonDeviceNotExist         = 101001 //设备不存在
	CodeCommonDeviceNotSupportedFunc = 101003 //设备不支持该功能
	CodeCommonFileIsEmpty            = 102001 // 空文件
	CodeCommonFileEncodingNotUTF8    = 102002 // 文件编码不是utf8

	//CodeFuncCollaboration
	CodeCollaborationOverLimit       = 100001
	CodeCollaborationDeviceOverLimit = 100002
	CodeCollaborationUserOverLimit   = 100003

	CodeCollaborationDismissed      = 100001
	CodeCollaborationUserNotExist   = 100002
	CodeCollaborationDeviceNotExist = 100003
	CodeCollaborationUrlExpired     = 100004
	CodeCollaborationUserExists     = 100005

	//CodeCommand
	CodeCommandParameterError              = 100001 //参数错误
	CodeCommandParamIndex                  = 100002 //命令类型错误
	CodeCommandParamTerminal               = 100003 //未选择终端
	CodeCommandParamGateway                = 100004 //未选择网关
	CodeCommandNotExist                    = 100005 //命令不存在
	CodeCommandCanNotCancel                = 100006 //命令不能取消
	CodeCommandCancelMultiLimit            = 100007 //批量操作最多支持1000个命令
	CodeCommandTimeoutInvalid              = 100008 //执行超时时间必须大于0
	CodeCommandSupportWhiteList            = 104001 //计划脱落命令只支持白名单模式
	CodeCommandExecTimeLessThanInvalidTime = 104002 //命令执行时间小于有效期时间
	CodeCommandScheduledArg2Invalid        = 104003 //命令计划执行时间不能小于当前时间

	//CodeFactorySetting
	CodeFactorySettingParameterError = 100001 //参数错误

	//网关配置
	CodeGatewaySettingModInvalid            = 100001 //无效的网关模式
	CodeGatewaySettingModNotSupport         = 100002 //设备固件不支持修改的网关模式
	CodeGatewaySettingInvalid               = 100003 // 网关无效，没有权限等情况
	CodeGatewaySettingMessageTypeInvalid    = 100004 // 网关消息类型无效
	CodeGatewaySettingTerminalDeviceInvalid = 100005 // 网关终端校验无效

	//file
	CodeFileEmpty      = 100001 //空文件
	CodeFileNotSupport = 100002 //只支持csv和xlsx文件

	//device
	CodeDeviceRecycled                   = 100001 //设备已被回收
	CodeDeviceFirmwareVersionInvalid     = 101001 //无效的固件版本号
	CodeDeviceBuildTimeInvalid           = 101002 //无效的固件编译时间
	CodeDeviceCodeHashInvalid            = 101003 //无效的固件代码哈希
	CodeDeviceMarkInvalid                = 102001 //无效的mark值
	CodeDeviceMarkMaxValue               = 103002 //mark最大值99999
	CodeDeviceRegisterTimeInvalid        = 106100 //出厂时间格式错误
	CodeDeviceUpdatedAtInvalid           = 107100 //数据更新时间格式错误
	CodeDeviceBatteryVoltageInvalid      = 108100 //无效电压值
	CodeDeviceInventoryStatusUnknown     = 109000 //设备未知状态
	CodeDeviceInventoryStatusDruidStock  = 109001 //druid库存状态
	CodeDeviceInventoryStatusAgencyStock = 109002 //agency库存状态
	CodeDeviceInventoryStatusSuspend     = 109003 //暂停状态
	CodeDeviceInventoryStatusRecycle     = 109004 //设备回收状态
	CodeDeviceInventoryStatusStop        = 109005 //设备停止状态
	CodeDeviceInventoryStatusNotInUse    = 109006 //设备不是使用状态
	CodeDeviceESNInvalid                 = 110001 //无效的ESN值
	CodeDeviceESNEmpty                   = 110002 //ESN为空
	CodeDeviceESNDuplicated              = 110003 //批量导入的文件中ESN重复了
	CodeDeviceHadAssigned                = 110004 //设备已经分配给其他用户
	CodeDeviceOnlyBindOneSatellite       = 110005 //设备只能绑定一个卫星

	CodeDeviceWithoutArgosModules        = 111001 //设备没有argos模组
	CodeDeviceAlreadyBeenBind            = 111002 //设备已经被绑定了
	CodeDeviceInvalidArgosID             = 111003 //无效的argos id
	CodeDeviceUnBindArgos                = 111004 //设备没有绑定argos
	CodeDeviceArgosAlreadyBeenBind       = 111005 //argos已经被绑定了
	CodeDeviceArgosR3ForbidChangeArgosID = 111006 //不可变更 R3 设备的 Argos ID
	CodeTargetArgosNotFound              = 111007 // 目标账号的Argos不存在
	CodeArgosUnsupportedChangeAccount    = 111008 // Argos不支持换绑， 需要先绑定

	CodeDeviceProductNotFind = 112001 // 该设备是未知型号

	//产品型号
	CodeProductParameterError = 100001 //参数错误

	//边缘智能
	CodeTriggerPinInvalid                = 100001 //无效的pin
	CodeTriggerPolygonPointNumOverLimit  = 100002 // 多边形围栏点数超过限制
	CodeTriggerSettingNotExist           = 100003 // 配置不存在
	CodeTriggerVHFDeviceRuleNumOverLimit = 100004 // 超过VHF设备规则数量限制

	//用户
	CodeUserParameterError    = 100001 //参数错误
	CodeUserNameInvalid_      = 100002 //无效用户id， CodeUserNameInvalid已被占用，CodeUserNameInvalid_表示
	CodeUserFullNameInvalid   = 100003 //无效用户名
	CodeUserPasswordInvalid   = 100004 //无效密码
	CodeUserEmailInvalid      = 100005 //无效邮箱
	CodeUserPhoneInvalid      = 100006 //无效电话号码
	CodeUserNameExist         = 100007 //用户ID已存在
	CodeUserPermissionInvalid = 100008 //用户无权限

	//配置
	CodeSettingParameterError                     = 100001 //参数错误
	CodeSettingDeviceNotSupportComposite          = 101008 // 设备不支持复合通信模式
	CodeSettingDeviceNotSupportBefore             = 101009 // 设备不支持跟随通信模式
	CodeSettingDeviceNotSupportSamplingContinuous = 101010 // 设备不支持持续采集
	CodeSettingDeviceNotSupportGprsPowerSaving    = 101011 // 设备不支持系统节能
	//固件管理
	CodeFirmwareUploadNameError     = 100001 //上传固件名字错误
	CodeFirmwareUploadFileError     = 100002 //上传固件文件错误
	CodeFirmwareUploadSizeError     = 100003 //上传固件文件大小错误
	CodeFirmwareUploadCodeHashError = 100004 //上传固件的代码哈希格式错误
	CodeFirmwareUploadExistError    = 100005 //上传固件的已存在
	// location
	CodeLocationFindError = 100001 // 定位数据查找失败

	// 功能配置文件
	CodeFeatureSettingFileDisabled = 100001 // 配置功能文件已禁用

	// Geo GeoIP
	CodeGetRealIpError     = 100001 // 获取真实IP失败
	CodeGetRealIpIsPrivate = 100002 // 获取真实IP为内网IP

	// Biological
	CodeBiologicalImageSizeExceedLimit = 103001 // 生物图片超过限制

	// VHF配置
	CodeVHFSettingInvalidPhyFrequency              = 100001 // 无效的中心频率
	CodeVHFSettingInvalidVoltageThreshold          = 100002 // 无效的电压门限
	CodeVHFSettingBroadcastSettingEmpty            = 100003 // 广播配置为空
	CodeVHFSettingBroadcastSettingInvalidMode      = 100004 // 无效的广播模式
	CodeVHFSettingBroadcastSettingInvalidInterval  = 100005 // 无效的广播间隔
	CodeVHFSettingBroadcastSettingInvalidDutyCycle = 100006 // 无效的广播时长
	CodeVHFSettingBroadcastSettingInvalidFrequency = 100007 // 无效的广播包间隔
	CodeVHFSettingBroadcastSettingInvalidPower     = 100008 // 无效的广播功率

	// 图片配置
	CodeCaptureImageSettingIntervalInvalid     = 100001 // 无效的图片采集间隔
	CodeCaptureImageSettingDurationTimeInvalid = 100002 // 无效的图片持续采集时间
	CodeCaptureImageSettingTimeTableIsZero     = 100003 // 图片采集时间点为0
	CodeCaptureImageSettingTimeTableIsInvalid  = 100004 // 无效的图片采集时间点
	CodeCaptureImageSettingModeInvalid         = 100005 // 无效的图片采集模式

	// 数据导出
	CodeExportTaskCancelError = 100001 // 数据导出任务不能取消
	CodeExportDeviceOverLimit = 100002 // 导出设备数超过限制

	//注册任务 CodeRegisterTask       = 50023
	CodeRegisterTaskNotExist         = 100001 // 注册任务不存在
	CodeRegisterTaskEXPIRE           = 100002 // 注册任务已过期
	CodeRegisterTaskRegistered       = 100003 // 注册任务已注册
	CodeRegisterTaskDeviceNotExist   = 101001 // 注册设备不存在
	CodeRegisterTaskDeviceRegistered = 101002 // 注册设备已注册

	// third device
	CodeThirdDeviceDeviceIDInvalid    = 100001 // 无效的三方设备id
	CodeThirdDeviceNotExist           = 100002 // 三方设备不存在
	CodeThirdDeviceFactoryInvalid     = 100003 // 无效的三方设备厂家
	CodeThirdDeviceModInvalid         = 100004 // 无效的三方设备产品型号
	CodeThirdDeviceDescriptionInvalid = 100005 // 无效的三方设备备注
	CodeThirdDeviceNotBelongToUser    = 100006 // 三方设备不属于用户
	CodeThirdDeviceImportNumLarge     = 100007 // 导入的三方设备过多
	// third device gps
	CodeThirdDeviceGPSImportNumLarge   = 105001 // 导入的三方设备GPS数量过多
	CodeThirdDeviceGPSLongitudeInvalid = 105002 // 导入的GPS经度不合法
	CodeThirdDeviceGPSLatitudeInvalid  = 105003 // 导入的GPS纬度不合法
	CodeThirdDeviceGPSAltitudeInvalid  = 105004 // 导入的GPS海拔高度不合法
	CodeThirdDeviceGPSCourseInvalid    = 105005 // 导入的GPS航向不合法
	CodeThirdDeviceGPSSpeedInvalid     = 105006 // 导入的GPS速度不合法
	CodeThirdDeviceGPSSatelliteInvalid = 105007 // 导入的GPS定位卫星数不合法
	CodeThirdDeviceGPSHDOPInvalid      = 105008 // 导入的GPS HDOP不合法
	CodeThirdDeviceGPSVDOPInvalid      = 105009 // 导入的GPS VDOP不合法
	CodeThirdDeviceGPSTimestampInvalid = 105010 // 导入的GPS 采集时间不合法
	CodeThirdDeviceGPSImportOtherData  = 105011 // 导入其他设备GPS数据
	// CodeEntrust 委托
	CodeEntrustBatchSupportDeviceNum = 100001 //批量委托最多支持1000个设备
	CodeEntrustNotExist              = 100002 // 委托不存在
	CodeEntrustNoPermission          = 100003 // 用户没权限操作此委托
	// 按文档修改
	CodeFirmwareLowerVersion = 102002 //
)

var (
	ErrDataSmall         = errors.New("Data too small.")
	ErrDataAlreadyExists = errors.New("Data already exists.")
)

type ErrorBody struct {
	Code    int    `json:"code"`
	SubCode []int  `json:"sub_code,omitempty"`
	Msg     string `json:"msg,omitempty"`
}

func (eb *ErrorBody) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", eb.Code, eb.Msg)
}

func IsErrBody(err interface{}) (errBody *ErrorBody, check bool) {
	if val, ok := err.(*ErrorBody); ok {
		return val, true
	} else if val, ok := err.(ErrorBody); ok {
		return &val, true
	}
	return nil, false
}

func HttpJsonErr(err interface{}) (int, *ErrorBody) {
	var (
		defaultHttpCode = http.StatusBadRequest // http 400 错误码才会被前端识别
		defaultErrCode  = CodeParameterInvalid
	)
	if val, ok := err.(*ErrorBody); ok {
		return defaultHttpCode, val
	} else if val, ok := err.(ErrorBody); ok {
		return defaultHttpCode, &val
	}
	return http.StatusInternalServerError, &ErrorBody{ // 非 ErrBody error 类型，说明err是非预期的
		Code: defaultErrCode,
		Msg:  fmt.Sprintf("%v", err),
	}
}

func NewErrorBody(code int, msg ...interface{}) *ErrorBody {
	return &ErrorBody{
		Code: code,
		Msg:  fmt.Sprint(msg...),
	}
}

func NewErrorBodyWithSingleSub(code int, subCode int, msg ...interface{}) *ErrorBody {
	return &ErrorBody{
		Code:    code,
		SubCode: []int{subCode},
		Msg:     fmt.Sprint(msg...),
	}
}

func NewErrorBodyWithSub(code int, subCode []int, msg ...interface{}) *ErrorBody {
	return &ErrorBody{
		Code:    code,
		SubCode: subCode,
		Msg:     fmt.Sprint(msg...),
	}
}

func Error(v ...interface{}) error {
	return errors.New(fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) error {
	return errors.New(fmt.Sprintf(format, v...))
}
