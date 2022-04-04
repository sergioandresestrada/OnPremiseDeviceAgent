import React from 'react'
import { Link } from 'react-router-dom'
import {Collapse, Nav, Navbar, NavbarBrand, NavbarToggler, NavItem} from 'reactstrap'

interface IHeader{
    isOpen: boolean
}

class Header extends React.Component<{}, IHeader>{
    constructor(props: any){
        super(props)

        this.state = { isOpen: false}
        this.toggle = this.toggle.bind(this)
    }

    toggle(){
        this.setState({
            isOpen: !this.state.isOpen
        })
    }

    render() {
        return(
            <div>
                <Navbar
                    color='light'
                    expand='xl'
                    full
                    light
                >
                    <NavbarBrand style={{ color: "#0096D6" }}>On-Premise Device Agent</NavbarBrand>
                    <NavbarToggler onClick={this.toggle}/>
                    <Collapse isOpen={this.state.isOpen} navbar>
                        <Nav navbar>
                            <NavItem>
                                    <Link to="/" style={{ color: "#0096D6", textDecoration: "none", padding:"1rem"}}>New Message</Link>
                            </NavItem>
                            
                            <NavItem>
                                    <Link to="/deviceInfoList" style={{ color: "#0096D6", textDecoration: "none", padding:"1rem"}}>Device Information</Link>
                            </NavItem>
                        </Nav>
                    </Collapse>
                </Navbar>
            </div>
        )
    }

}

export default Header;