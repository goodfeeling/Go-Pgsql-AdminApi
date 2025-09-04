package task_execution_log

import (
	"encoding/json"
	"log"

	domainTaskExecutionLog "github.com/gbrayhan/microservices-go/src/domain/sys/task_execution_log"
	ws "github.com/gbrayhan/microservices-go/src/infrastructure/lib/websocket"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// LogHandler 日志处理器
type LogHandler struct {
	Logger                  *logger.Logger
	taskExecutionLogService domainTaskExecutionLog.ITaskExecutionLogService
	wsRouter                *ws.WebSocketRouter
	context                 *ws.WebSocketContext
}

func NewLogHandler(
	taskExecutionLogService domainTaskExecutionLog.ITaskExecutionLogService,
	loggerInstance *logger.Logger,
	wsRouter *ws.WebSocketRouter,
) *LogHandler {
	return &LogHandler{
		Logger:                  loggerInstance,
		taskExecutionLogService: taskExecutionLogService,
	}
}

// 实现扩展接口
func (ch *LogHandler) OnConnectWithContext(conn *websocket.Conn, ctx *ws.WebSocketContext) {
	ch.Logger.Info("Log handler: Client connected")
}

func (ch *LogHandler) OnConnect(conn *websocket.Conn) {
	ch.Logger.Info("Log handler: Client connected")

	result, err := ch.taskExecutionLogService.GetByTaskID(0, 100)
	if err != nil {
		ch.sendError(conn, "Failed to fetch initial logs: "+err.Error())
		return
	}
	// 发送分页结果
	ch.sendData(conn, "initial_logs", result)
}

func (ch *LogHandler) OnMessage(conn *websocket.Conn, message []byte) {
	ch.Logger.Info("Log handler: Received message: %s", zap.String("Message", string(message)))
	// 解析传入的 JSON 字符串
	var request struct {
		TaskID int `json:"taskId"`
		Limit  int `json:"limit"`
	}
	if err := json.Unmarshal(message, &request); err != nil {
		ch.sendError(conn, "Invalid JSON format: "+err.Error())
		return
	}

	ch.Logger.Info("Log handler: Fetching logs for task ID:", zap.Int("TaskID", request.TaskID))

	// 设置默认值
	if request.Limit <= 0 {
		request.Limit = 100
	}

	// 调用 GetByTaskID 获取日志数据
	result, err := ch.taskExecutionLogService.GetByTaskID(uint(request.TaskID), request.Limit)
	if err != nil {
		ch.sendError(conn, "Failed to fetch logs: "+err.Error())
		return
	}
	ch.Logger.Info("Log handler: Fetched logs: %v", zap.Any("Logs", len(*result)))

	// 发送结果
	ch.sendData(conn, "logs", result)
}

// 当单条数据发生变化时调用此方法
func (ch *LogHandler) NotifySingleDataChange(action string, data interface{}) {
	// action 可以是 "create", "update", "delete"
	response := map[string]interface{}{
		"type":   "data_changed",
		"action": action,
		"data":   data,
	}

	jsonData, _ := json.Marshal(response)
	ch.wsRouter.WebSocketManager.BroadcastMessage(jsonData)
}

// 当多条数据发生变化时调用此方法
func (ch *LogHandler) NotifyBatchDataChange(action string, dataList []interface{}) {
	response := map[string]interface{}{
		"type":   "batch_data_changed",
		"action": action,
		"data":   dataList,
	}

	jsonData, _ := json.Marshal(response)
	ch.wsRouter.WebSocketManager.BroadcastMessage(jsonData)
}

func (ch *LogHandler) OnDisconnect(conn *websocket.Conn) {
	log.Println("Chat handler: Client disconnected")
}
func (ch *LogHandler) OnDisconnectWithContext(conn *websocket.Conn, ctx *ws.WebSocketContext) {
	ch.Logger.Info("Log handler: Client disconnected")
	// 连接断开时清理会话

}

// 辅助方法
func (ch *LogHandler) sendData(conn *websocket.Conn, msgType string, data interface{}) {
	response := map[string]interface{}{
		"type": msgType,
		"data": data,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		ch.sendError(conn, "Failed to marshal response")
		return
	}

	conn.WriteMessage(websocket.TextMessage, jsonData)
}

func (ch *LogHandler) sendError(conn *websocket.Conn, message string) {
	errorResponse := map[string]interface{}{
		"type":  "error",
		"error": message,
	}

	jsonData, _ := json.Marshal(errorResponse)
	conn.WriteMessage(websocket.TextMessage, jsonData)
}
