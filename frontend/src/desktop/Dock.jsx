import React from 'react';
import '../styles/dock.css';

export default function Dock({ onOpenWindow, apps }) {
  return (
    <div className="dock-container">
      <div className="dock">
        {apps && apps.map(app => (
          <div className="dock-item" key={app.id} onClick={() => onOpenWindow(app.window)}>
            <img src={app.icon} alt={app.name} />
            <span className="tooltip">{app.name}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
