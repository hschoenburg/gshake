import React, { Component } from 'react';

class NotifyForm extends Component {

  constructor(props) { 
    super(props)
    this.props = props
    this.Notify = this.Notify.bind(this)
  }

  render() {
    return (
      <div>
        <p>Want to be notified before "{this.props.name}" can be opened for bidding?</p>
        <form onSubmit={this.Notify}>
          <input 
            placeholder="Email"
            ref="email"
          />
          <input type = 'submit' />
        </form>
      </div>
    )
  }

  async Notify (e) {
    e.preventDefault()
    try {
      const url = process.env.REACT_APP_GSHAKE_HOST + "/notify"
      const postOpts = {
        method: 'POST',
        headers: { "ContentType": "application/json" },
        body: JSON.stringify({contact: this.refs.email.value, name: this.props.name })
      }
      console.log(postOpts)

      let req = await fetch(url, postOpts)
      let res = await req.json()
      console.log(res)
      //this.setState({info: info})
    } catch (err) {
      console.log('NoifyForm: Notify Error: ' + err)
    }
  }
}

export default NotifyForm;
