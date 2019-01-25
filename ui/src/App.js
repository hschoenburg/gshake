import React, { Component } from 'react';
import logo from './hns-logo.svg';
import './App.css';
import InfoBox from './InfoBox';

class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>HNSNNS</h2>
          <h4>Handshake Name Notifier Service</h4> 
          <InfoBox />
        </header>
      </div>
    );
  }
}

export default App;
