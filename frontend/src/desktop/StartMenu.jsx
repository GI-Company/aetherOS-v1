import React, { useState } from 'react';
import '../styles/startMenu.css';

const apps = [
  { id: 'ai', name: 'AI Agent', icon: '/aetherOS_icon.png' },
  { id: 'compute', name: 'Compute', icon: '/compute-icon.png' },
  { id: 'marketplace', name: 'Marketplace', icon: '/marketplace-icon.png' },
];

export default function StartMenu({ onOpenWindow }) {
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  const toggleMenu = () => setIsOpen(!isOpen);

  const handleDragStart = (e, appId) => {
    e.dataTransfer.setData('appId', appId);
  };

  const filteredApps = apps.filter(app =>
    app.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="start-menu-container">
      <button className="start-button" onClick={toggleMenu}>
        <img src="/aetherOS_icon.png" alt="Start" />
      </button>
      {isOpen && (
        <div className="start-menu">
          <div className="search-bar">
            <input
              type="text"
              placeholder="Search apps..."
              value={searchTerm}
              onChange={e => setSearchTerm(e.target.value)}
            />
          </div>
          <div className="app-list">
            {filteredApps.map(app => (
              <div
                key={app.id}
                className="app-list-item"
                onClick={() => { onOpenWindow(app.id); toggleMenu(); }}
                draggable
                onDragStart={e => handleDragStart(e, app.id)}
              >
                <img src={app.icon} alt={app.name} />
                <span>{app.name}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
