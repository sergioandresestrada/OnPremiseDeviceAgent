import React from 'react'
import { Link } from 'react-router-dom'
import {Collapse, DropdownItem, DropdownMenu, DropdownToggle, Nav, Navbar, NavbarBrand, NavbarToggler, NavItem, UncontrolledDropdown} from 'reactstrap'

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
                    <NavbarBrand style={{ color: "#0096D6", marginRight: "2em"}}>On-Premise Device Agent</NavbarBrand>
                    <NavbarToggler onClick={this.toggle}/>
                    <Collapse isOpen={this.state.isOpen} navbar>
                        <Nav navbar>
                            <NavItem style={{ padding:"0.5rem"}}>
                                    <Link to="/" style={{ color: "#0096D6", textDecoration: "none"}}>New Message</Link>
                            </NavItem>
                            
                            <NavItem style={{ color: "#0096D6", padding:"0.5rem", marginLeft: "0.5em"}}>
                                    <Link to="/deviceInfoList" style={{ color: "#0096D6", textDecoration: "none"}}>Requested Device Information</Link>
                            </NavItem>
                            <UncontrolledDropdown inNavbar nav style={{ marginLeft: "0.5em"}}>
                                <DropdownToggle caret nav style={{ color: "#0096D6", textDecoration: "none", padding:"0.5rem"}}>
                                    Devices
                                </DropdownToggle>
                                <DropdownMenu right>
                                    <DropdownItem>
                                        <Link to="/devices/new" style={{ color: "#0096D6", textDecoration: "none"}}>New</Link>
                                    </DropdownItem>
                                    <DropdownItem>
                                        <Link to="/devices" style={{ color: "#0096D6", textDecoration: "none"}}>List</Link>
                                    </DropdownItem>
                                </DropdownMenu>
                            </UncontrolledDropdown>
                            
                        </Nav>
                    </Collapse>
                </Navbar>
            </div>
        )
    }

}

export default Header;