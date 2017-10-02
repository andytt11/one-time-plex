import React, { Component } from 'react'
import logo from './logo.svg'
import Users from './containers/users'
import AddUser from './containers/adduser'
import Search from './containers/search'
import { Grid, Cell } from 'react-mdl'
import './App.css'

window.otp = {
  url: 'http://localhost:8080/api'
}

class App extends Component {
  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Welcome to One Time Plex</h2>
        </div>

        <div style={{ width: '80%', margin: 'auto'}}>
          <Grid>
            <Cell col={12}>
              <h4>Restricted Users:</h4>

              <Users />
            </Cell>
            <Cell col={6}>
              <h4>Add User:</h4>
              
              <AddUser />
            </Cell>

            <Cell col={6}>
              <h4>Search Plex:</h4>

              <Search />
            </Cell>
          </Grid>
        </div>
      </div>
    )
  }
}

export default App