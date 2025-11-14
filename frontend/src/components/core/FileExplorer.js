import React, { useState, useEffect } from 'react';
import { send, on, off } from '../../services/websocket';

function FileExplorer() {
  const [files, setFiles] = useState([]);
  const [currentPath, setCurrentPath] = useState('/');

  useEffect(() => {
    const handleVfsListResult = (payload) => {
      setFiles(payload.files);
    };

    on('vfs:list:result', handleVfsListResult);
    send('vfs:list', { path: currentPath });

    return () => {
      off('vfs:list:result', handleVfsListResult);
    };
  }, [currentPath]);

  const handleFileClick = (file) => {
    if (file.type === 'folder') {
      const newPath = currentPath === '/' ? `/${file.name}` : `${currentPath}/${file.name}`;
      setCurrentPath(newPath);
    }
  };

  const handleGoBack = () => {
    if (currentPath !== '/') {
      const newPath = currentPath.substring(0, currentPath.lastIndexOf('/')) || '/';
      setCurrentPath(newPath);
    }
  };

  const handleNewFile = () => {
    const fileName = prompt('Enter file name:');
    if (fileName) {
      send('vfs:create:file', { path: `${currentPath}/${fileName}` });
    }
  };

  const handleNewFolder = () => {
    const folderName = prompt('Enter folder name:');
    if (folderName) {
      send('vfs:create:folder', { path: `${currentPath}/${folderName}` });
    }
  };

  const handleDelete = (file) => {
    if (window.confirm(`Are you sure you want to delete ${file.name}?`)) {
      send('vfs:delete', { path: `${currentPath}/${file.name}` });
    }
  };

  const renderFiles = () => {
    return files.map((file, index) => (
      <li key={index}>
        <span onClick={() => handleFileClick(file)} style={{ cursor: 'pointer' }}>
          {file.type === 'folder' ? '📁' : '📄'}
          {file.name}
        </span>
        <button onClick={() => handleDelete(file)} style={{ marginLeft: '10px' }}>Delete</button>
      </li>
    ));
  };

  return (
    <div>
      <h2>File Explorer</h2>
      <p>Current Path: {currentPath}</p>
      {currentPath !== '/' && <button onClick={handleGoBack}>Go Back</button>}
      <button onClick={handleNewFile}>New File</button>
      <button onClick={handleNewFolder}>New Folder</button>
      <ul>
        {renderFiles()}
      </ul>
    </div>
  );
}

export default FileExplorer;
