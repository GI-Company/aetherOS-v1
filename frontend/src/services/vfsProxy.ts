import { openDB, IDBPDatabase } from 'idb';

class VFSProxy {
  private ws: WebSocket | null = null;
  private db: IDBPDatabase | null = null;
  private reconnectAttempts = 0;

  constructor(private wsUrl: string) {
    this.initDB();
    this.connect();
  }

  private async initDB() {
    this.db = await openDB('AetherVFS', 1, {
      upgrade(db) {
        db.createObjectStore('files');
      },
    });
  }

  private connect() {
    this.ws = new WebSocket(this.wsUrl);

    this.ws.onopen = () => {
      console.log('VFSProxy connected to kernel');
      this.reconnectAttempts = 0;
    };

    this.ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      this.handleKernelMessage(msg);
    };

    this.ws.onclose = () => {
      console.log('VFSProxy disconnected from kernel');
      this.reconnect();
    };

    this.ws.onerror = (error) => {
      console.error('VFSProxy WebSocket error:', error);
      // The onclose event will be called next, which will trigger the reconnect logic.
    };
  }

  private reconnect() {
    if (this.reconnectAttempts >= 5) {
      console.error('VFSProxy: Max reconnection attempts reached. Giving up.');
      return;
    }

    this.reconnectAttempts++;
    const delay = Math.pow(2, this.reconnectAttempts) * 1000;
    console.log(`VFSProxy: Reconnecting in ${delay / 1000} seconds...`);

    setTimeout(() => {
      this.connect();
    }, delay);
  }

  private async handleKernelMessage(msg: any) {
    if (!this.db) return;

    const { topic, payload } = msg;

    switch (topic) {
      case 'vfs:write':
        try {
          await this.db.put('files', payload.content, payload.path);
          this.sendAck(msg.id, 'success');
        } catch (error) {
          console.error('VFSProxy: Error writing to IndexedDB:', error);
          this.sendAck(msg.id, 'error');
        }
        break;
      case 'vfs:read':
        try {
          const content = await this.db.get('files', payload.path);
          this.sendMessage('vfs:read:resp', { path: payload.path, content });
        } catch (error) {
          console.error('VFSProxy: Error reading from IndexedDB:', error);
        }
        break;
      case 'vfs:list':
        try {
          const keys = await this.db.getAllKeys('files');
          this.sendMessage('vfs:list:resp', { keys });
        } catch (error) {
          console.error('VFSProxy: Error listing keys from IndexedDB:', error);
        }
        break;
    }
  }

  public sendMessage(topic: string, payload: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ topic, payload }));
    } else {
      console.error('VFSProxy: WebSocket is not open. Message not sent.');
    }
  }

  private sendAck(id: string, status: string) {
    this.sendMessage('ack', { id, status });
  }
}

export default VFSProxy;
