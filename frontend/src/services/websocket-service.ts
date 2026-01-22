// WebSocket service for real-time messaging

const WS_BASE_URL = import.meta.env.VITE_API_BASE_URL?.replace('http', 'ws') || 'ws://localhost:8082';

export type WebSocketMessageType =
  | 'new_message'
  | 'send_message'
  | 'ack_message'
  | 'typing'
  | 'in_chat'
  | 'new_notification'
  | 'error';

export interface WebSocketMessage {
  type: WebSocketMessageType;
  payload: any;
}

// Handler nên nhận payload trực tiếp để code sạch hơn
type MessageHandler = (payload: any) => void;

class WebSocketService {
  private ws: WebSocket | null = null;
  private handlers: Map<WebSocketMessageType, MessageHandler[]> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private isConnecting = false;
  private shouldReconnect = true;

  connect(token: string): Promise<void> {
    if (this.ws?.readyState === WebSocket.OPEN) return Promise.resolve();
    if (this.isConnecting) return Promise.resolve();

    return new Promise((resolve, reject) => {
      try {
        this.isConnecting = true;
        // Đảm bảo URL này khớp với RegisterWebSocketRoutes ở BE
        const wsUrl = `${WS_BASE_URL}/api/ws?token=${token}`;
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
          console.log('✅ WebSocket Connected');
          this.isConnecting = false;
          this.reconnectAttempts = 0;
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const message: WebSocketMessage = JSON.parse(event.data);
            // TRUYỀN PAYLOAD VÀO ĐÂY
            this.handleMessage(message.type, message.payload);
          } catch (error) {
            console.error('❌ Failed to parse WS message:', error);
          }
        };

        this.ws.onerror = (error) => {
          this.isConnecting = false;
          reject(error);
        };

        this.ws.onclose = () => {
          this.isConnecting = false;
          this.ws = null;
          if (this.shouldReconnect && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            setTimeout(() => this.connect(token), this.reconnectDelay);
            this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
          }
        };
      } catch (error) {
        this.isConnecting = false;
        reject(error);
      }
    });
  }

  send(type: WebSocketMessageType, payload: any): boolean {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return false;
    try {
      this.ws.send(JSON.stringify({ type, payload }));
      return true;
    } catch (error) {
      return false;
    }
  }

  on(type: WebSocketMessageType, handler: MessageHandler) {
    if (!this.handlers.has(type)) this.handlers.set(type, []);
    this.handlers.get(type)!.push(handler);
  }

  off(type: WebSocketMessageType, handler: MessageHandler) {
    const handlers = this.handlers.get(type);
    if (handlers) {
      const index = handlers.indexOf(handler);
      if (index > -1) handlers.splice(index, 1);
    }
  }

  private handleMessage(type: WebSocketMessageType, payload: any) {
    const handlers = this.handlers.get(type);
    if (handlers) {
      handlers.forEach(handler => handler(payload)); // Chạy với payload
    }
  }

  disconnect() {
    this.shouldReconnect = false;
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }
}

export const wsService = new WebSocketService();