import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Login from './components/auth/Login';
import HandleAuth from './components/auth/HandleAuth';
import FileExplorer from './components/core/FileExplorer';

function App() {
  return (
    <Router>
      <div>
        <nav>
          <ul>
            <li>
              <Link to="/">Home</Link>
            </li>
            <li>
              <Link to="/login">Login</Link>
            </li>
          </ul>
        </nav>

        <hr />

        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/auth" element={<HandleAuth />} />
          <Route path="/" element={<FileExplorer />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
