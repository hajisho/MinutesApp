"use strict"
import React, { useEffect, useState } from "react"
import ReactDOM from 'react-dom'

function GetMessage() {
  const [data, setData] = useState({ hits: [] });

  useEffect(() => {

    const data = async() => {
      const res = await fetch("/message")
                                .then(response => response.json())
                                .then(data => console.log(data[0].message))
    }
    data()

  }, []);

  return (
    <div>
      <h1>Cool app</h1>
    </div>
  );
}

ReactDOM.render(<GetMessage />, document.getElementById('message'));


class Hello extends React.Component{
  render(){
    return(
      <div>Hello World</div>
    );
  }
}

ReactDOM.render(<Hello />, document.getElementById("content"))
