import React, { useState, useEffect } from 'react';

export default function Taskbar({ windows }) {
  const [time, setTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => setTime(new Date()), 1000);
    return () => clearInterval(timer);
  }, []);

  return (
    <div className="taskbar">
      <div className="taskbar-windows">
        {windows.map((w) => (
          <button key={w.id} onClick={w.toggleMinimize}>
            {w.title}
          </button>
        ))}
      </div>
      <div className="taskbar-clock">
        {time.toLocaleTimeString()}
      </div>
    </div>
  );
}