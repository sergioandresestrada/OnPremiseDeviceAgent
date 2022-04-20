import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import reportWebVitals from './reportWebVitals';
import 'bootstrap/dist/css/bootstrap.min.css';
import Header from './Components/Header';
import Form from './Components/Form';
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import DeviceInfoList from './Components/DeviceInfoList';
import RequestedDeviceInfo from './Components/RequestedDeviceInfo';
import DeviceForm from './Components/DeviceForm';
import DeviceList from './Components/DeviceList';

ReactDOM.render(
  <React.StrictMode>
    <Router>
      <header>
        <Header />
      </header>
      <Routes>
        <Route path="/" element={<Form />} />
        <Route path="/deviceInfoList" element={<DeviceInfoList/>} />
        <Route path="/deviceInfo" element={<RequestedDeviceInfo/>} />
        <Route path="/devices/new" element={<DeviceForm {...{isNewDevice: true}}/>} />
        <Route path="/devices" element={<DeviceList/>} />
        <Route path="/devices/edit/:uuid" element={<DeviceForm {...{isNewDevice: false}} />} />
        <Route
          path="*"
          element={
            <main style={{ padding: "1rem" }}>
              <p>There's nothing here!</p>
            </main>
          }
        />
      </Routes>
    </Router>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
