import React, { Component } from 'react';
import './InfoBox.css';
import NotifyForm from './NotifyForm'

class InfoBox extends Component {

  state = { info: {}, error: false }

  render() {
    return (
      <div className="info-box">
        <p>Look up Name</p>
        <form onSubmit={this.getInfo}>
          <input 
            placeholder="Name"
            ref="name"
          />

          <input type = 'submit' />
        </form>
        <SearchResults {...this.state}/>
      </div>
    );
  }

  getInfo = async (e) => {
    e.preventDefault()
    const name = this.refs.name.value
    if(name.length === 0) { return }
    let req = await fetch(process.env.REACT_APP_GSHAKE_HOST + "/info/" + name)
    this.refs.name.value = ''
    let res = await req.json()
    if(req.ok) {
      this.setState({info: res, error: false})
    } else {
      this.setState({info:{} , error: {message: res, name: name}})
    }
  }
}

const SearchResults = (props) => {

  if(props.error.message) {
    return (
      <div>
          <p>Sorry an Error Occured. Please try again</p>
          <p>{props.name} Error:  {props.error.message}</p>
      </div>
    ) 
  } else if(props.info.Name) {
    let info = props.info
    return (
      <div>
        <ul className='search-results'>
          <li>Name: {info.Name}</li>
          <li>Reserved: {String(info.Result.start.reserved)}</li>
          <li>Week Available: {info.Result.start.week}</li>
          <li>Block Height: {info.Result.start.start}</li>
        </ul>
        <NotifyForm name={info.Name} />
      </div>
    )
  } else {
    return null
  }
}


export default InfoBox;
