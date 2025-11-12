import React from 'react';

export default function Taskbar({ windows }) {
  return (
    <div className="taskbar">
      {windows.map((w) => (
        <button key={w.id} onClick={w.toggleMinimize}>
          {w.title}
        </button>
      ))}
    </div>
  );
}