import React from 'react';
import Login from './components/Login/Login';
import Signup from './components/Signup/Signup';
import Home from './pages/Home';
import {BrowserRouter, Route, Switch } from 'react-router-dom';
import './App.scss';
import { AuthServiceClient } from './proto/services_grpc_web_pb';

export const authClient = new AuthServiceClient("http://localhost:9001")

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Switch>
          <Route path="/" exact component={Home} />
        </Switch>
        <Switch>
          <Route path="/login" exact component={Login} />
        </Switch>
        <Switch>
          <Route path="/signup" exact component={Signup} />
        </Switch>
      </BrowserRouter>
    </div>
  );
}

export default App;
