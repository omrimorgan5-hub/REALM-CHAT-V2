import React, { useEffect, useState, useRef } from 'react';
import './App.css'; // Import your new CSS!

interface ChatMessage {
  username: string;
  content: string;
  timestamp: string;
}

const App: React.FC = () => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [username] = useState(() => {
  const savedName = localStorage.getItem('realm_username');
  
  if (savedName) {
    return savedName; // Use the name they had before
  } else {
    // 2. If it's their first time, make a new one and save it
    const newName = `Anon-${Math.floor(Math.random() * 9999)}`;
    localStorage.setItem('realm_username', newName);
    return newName;
  }
});
  const socketRef = useRef<WebSocket | null>(null);
  const chatEndRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    chatEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  useEffect(() => {
    // Using 127.0.0.1 to avoid Firefox/Windows localhost issues
    const socket = new WebSocket('ws://192.168.1.112:8080/ws');
    socketRef.current = socket;

    socket.onopen = () => console.log("Connected to the Realm");
    
    socket.onmessage = (event) => {
      const data: ChatMessage = JSON.parse(event.data);
      setMessages((prev) => [...prev, data]);
    };

    socket.onerror = (err) => console.error("Socket Error:", err);

    return () => socket.close();
  }, []);

  const sendMessage = () => {
    if (input.trim() && socketRef.current?.readyState === WebSocket.OPEN) {
      const msg: ChatMessage = {
        username,
        content: input,
        timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      };
      socketRef.current.send(JSON.stringify(msg));
      setInput('');
    }
  };

  return (
    <div className="realm-container">
      <header className="realm-header">
        <h1>REALM</h1>
        <p>Anonymous Entry: {username}</p>
      </header>

      <div className="chat-box">
        {messages.map((msg, i) => (
          <div key={i} className="message">
            <span className="msg-time">{msg.timestamp}</span>
            <span className="msg-user">{msg.username}</span>
            <span className="msg-content">{msg.content}</span>
          </div>
        ))}
        <div ref={chatEndRef} />
      </div>

      <div className="input-area">
        <input 
          value={input} 
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
          placeholder="Transmit to the Realm..."
        />
        <button onClick={sendMessage}>Send</button>
      </div>
    </div>
  );
};

export default App;