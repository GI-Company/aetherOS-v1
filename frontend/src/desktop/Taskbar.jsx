import React, { useState, useEffect } from 'react';
import StartMenu from './StartMenu';
import '../styles/taskbar.css';

export default function Taskbar({ windows, onOpenWindow }) {
  const [dateTime, setDateTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => setDateTime(new Date()), 1000);
    return () => clearInterval(timer);
  }, []);

  return (
    <div className="taskbar">
      <StartMenu onOpenWindow={onOpenWindow} />
      <div className="taskbar-windows">
        {windows.map((w) => (
          <button key={w.id} onClick={w.toggleMinimize}>
            {w.title}
          </button>
        ))}
      </div>
      <div className="taskbar-clock">
        <span>{dateTime.toLocaleDateString()}</span>
        <span>{dateTime.toLocaleTimeString()}</span>
      </div>
    </div>
  );
}
