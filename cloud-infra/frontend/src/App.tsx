import React from 'react';
import Form from './Form';
import Header from './Header';
import './App.css';

function App() {
  return (
    <div className="App">
      <header>
        <Header />
      </header>
      <div className='Form'>
        <Form/>
      </div>
    </div>     

  );
}

export default App;
