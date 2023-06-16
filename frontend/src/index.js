import React from 'react'
import ReactDOM from 'react-dom'
import {
  BrowserRouter as Router,
  Route,
  Switch,
  Redirect
} from "react-router-dom";
import NoAuthRoute from "./component/NoAuthRoute"
import AuthRoute from "./component/AuthRoute"
import AddNewItem from "./component/AddNewItem"
import NewUser from "./component/NewUser"
import SignIn from "./component/SignIn"
import SignUp from "./component/SignUp"
import SignOutRoute from "./component/SignOutRoute"
import NotYou from "./component/NotYou"
import Verify from "./component/Verify"
import VerifyLink from "./component/VerifyLink"
import Home from "./component/Home"
import Search from "./component/Search"
import Followers from "./component/Followers"
import ErrorPage from "./component/ErrorPage"
// import Cookies from 'js-cookie';
import api from './api'
import 'bootstrap/dist/css/bootstrap.min.css'
import './assets/stylesheet.css'
import { combineReducers, configureStore } from '@reduxjs/toolkit';
import userReducer from "./reducers/user";
import authReducer from "./reducers/auth";
import { Provider } from 'react-redux';
import storage from 'redux-persist/lib/storage';
import { persistReducer, persistStore } from 'redux-persist';
import thunk from 'redux-thunk';
import { PersistGate } from 'redux-persist/integration/react';

const persistConfig = {
  key: 'root',
  storage,
}

const rootReducer = combineReducers({ 
  user: userReducer,
  auth: authReducer
})

const persistedReducer = persistReducer(persistConfig, rootReducer)

const store = configureStore({
  reducer: persistedReducer,
  middleware: [thunk]
});

const persistor = persistStore(store)

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <PersistGate loading={null} persistor={persistor}>
        <Router>
          <Switch>
            <NoAuthRoute path="/sign-up" component={SignUp}/>
            <NoAuthRoute path="/sign-in" component={SignIn}/>
            <Route path="/verify" component={Verify}/>
            <Route path="/verify-link" component={VerifyLink}/>
            <Route path="/not-you" component={NotYou}/>
            <Route path="/oh-no" component={ErrorPage}/>
            <AuthRoute path="/new-user" component={NewUser}/>
            <AuthRoute exact path="/" redirect="/sign-up" component={Home}/>
            <AuthRoute path="/search" component={Search}/>
            <AuthRoute path="/followers" component={Followers}/>
            <AuthRoute path="/add-new" component={AddNewItem}/>
            <SignOutRoute path="/sign-out"/>
            <Redirect to={{pathname: "/oh-no", state: { errorCode: "not_found"}}}/>
          </Switch>
        </Router>
      </PersistGate>
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
)
