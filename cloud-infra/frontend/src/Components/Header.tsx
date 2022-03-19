import React from 'react'
import {Collapse, Nav, Navbar, NavbarBrand, NavbarToggler, NavItem, NavLink} from 'reactstrap'

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
                                <NavLink style={{ color: "#0096D6" }}>New Message</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink style={{ color: "#0096D6" }}>Placeholder</NavLink>
                            </NavItem>
                        </Nav>
                    </Collapse>
                </Navbar>
            </div>
        )
    }

}

export default Header;