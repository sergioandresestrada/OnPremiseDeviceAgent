import React from 'react';
import Form from './Components/Form';
import Header from './Components/Header';
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
