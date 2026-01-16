// WebSocket service for real-time messaging

const WS_BASE_URL = import.meta.env.VITE_API_BASE_URL?.replace('http', 'ws') || 'ws://localhost:8080';

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

export interface SendMessagePayload {
  channel_id: string;
  content: string;
  temp_id?: string; // For optimistic UI updates
}

export interface ACKMessagePayload {
  temp_message_id: string;
  message: {
    id: string;
    channel_id: string;
    sender_id: string;
    sender_username: string;
    type: 'user' | 'system';
    content: string;
    read_by: string[];
    created_at: string;
  };
}

export interface IncomingMessagePayload {
  message: {
    id: string;
    channel_id: string;
    sender_id: string;
    sender_username: string;
    type: 'user' | 'system';
    content: string;
    read_by: string[];
    created_at: string;
  };
}

export interface ErrorMessagePayload {
  error_msg: string;
  temp_message_id?: string;
}

type MessageHandler = (message: WebSocketMessage) => void;

class WebSocketService {
  private ws: WebSocket | null = null;
  private handlers: Map<WebSocketMessageType, MessageHandler[]> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000; // Start with 1 second
  private isConnecting = false;
  private shouldReconnect = true;

  /**
   * Connect to WebSocket server
   */
  connect(token: string): Promise<void> {
    if (this.ws?.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected');
      return Promise.resolve();
    }

    if (this.isConnecting) {
      console.log('WebSocket connection in progress');
      return Promise.resolve();
    }

    return new Promise((resolve, reject) => {
      try {
        this.isConnecting = true;
        const wsUrl = `${WS_BASE_URL}/api/ws?token=${token}`;
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
          console.log('WebSocket connected');
          this.isConnecting = false;
          this.reconnectAttempts = 0;
          this.reconnectDelay = 1000;
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const message: WebSocketMessage = JSON.parse(event.data);
            this.handleMessage(message);
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
          }
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.isConnecting = false;
          reject(error);
        };

        this.ws.onclose = (event) => {
          console.log('WebSocket closed:', event.code, event.reason);
          this.isConnecting = false;
          this.ws = null;

          // Attempt reconnect if not intentionally closed
          if (this.shouldReconnect && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`Reconnecting... attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}`);
            
            setTimeout(() => {
              this.connect(token).catch(err => {
                console.error('Reconnect failed:', err);
              });
            }, this.reconnectDelay);

            // Exponential backoff
            this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
          }
        };
      } catch (error) {
        this.isConnecting = false;
        reject(error);
      }
    });
  }

  /**
   * Disconnect from WebSocket
   */
  disconnect() {
    this.shouldReconnect = false;
    if (this.ws) {
      this.ws.close(1000, 'Client disconnect');
      this.ws = null;
    }
  }

  /**
   * Send a message through WebSocket
   */
  send(type: WebSocketMessageType, payload: any): boolean {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not connected');
      return false;
    }

    try {
      const message: WebSocketMessage = { type, payload };
      this.ws.send(JSON.stringify(message));
      return true;
    } catch (error) {
      console.error('Failed to send WebSocket message:', error);
      return false;
    }
  }

  /**
   * Subscribe to specific message type
   */
  on(type: WebSocketMessageType, handler: MessageHandler) {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, []);
    }
    this.handlers.get(type)!.push(handler);
  }

  /**
   * Unsubscribe from message type
   */
  off(type: WebSocketMessageType, handler: MessageHandler) {
    const handlers = this.handlers.get(type);
    if (handlers) {
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    }
  }

  /**
   * Handle incoming message
   */
  private handleMessage(message: WebSocketMessage) {
    const handlers = this.handlers.get(message.type);
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(message);
        } catch (error) {
          console.error(`Error in handler for ${message.type}:`, error);
        }
      });
    }
  }

  /**
   * Check if WebSocket is connected
   */
  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }
}

// Export singleton instance
export const wsService = new WebSocketService();
