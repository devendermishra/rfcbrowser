import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import RFC from './RFC';
import RFCList from './RFCList';


const App: React.FC = () => (
  <Router>
    <Routes>
      <Route path="/rfc/:id" element={<RFC />} />
      <Route path="/" element={<RFCList />} />
    </Routes>
  </Router>
);

export default App;
