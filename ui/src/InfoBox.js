import React, { Component } from 'react';
import './InfoBox.css';
import NotifyForm from './NotifyForm'

class InfoBox extends Component {

  state = { info: {}}

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
        <SearchResults {...this.state.info}/>
      </div>
    );
  }

  getInfo = async (e) => {
      e.preventDefault()
      const name = this.refs.name.value
      let req = await fetch(process.env.REACT_APP_GSHAKE_HOST + "/info/" + name)
      this.refs.name.value = ''
      let info = await req.json()
      this.setState({info: info})
  }
}

const SearchResults = (props) => {
  if(props.Name) {
    return (
      <div>
        <ul className='search-results'>
          <li>Name: {props.Name}</li>
          <li>Reserved: {String(props.Result.start.reserved)}</li>
          <li>Week Available: {props.Result.start.week}</li>
          <li>Block Height: {props.Result.start.start}</li>
        </ul>
        <NotifyForm name={props.Name} />
      </div>
    )
  } else {
    return null
  }
}


export default InfoBox;
