import React from 'react'
import {Collapse, Nav, Navbar, NavbarBrand, NavItem, NavLink} from 'reactstrap'

class Header extends React.Component{
    render() {
        return(
            <div>
                <Navbar
                    color='light'
                    expand='xl'
                    full
                >
                    <NavbarBrand>On-Premise Device Agent</NavbarBrand>
                    <Collapse navbar>
                        <Nav navbar>
                            <NavItem>
                                <NavLink href="/">New Job</NavLink>
                            </NavItem>
                        </Nav>
                    </Collapse>
                </Navbar>
            </div>
        )
    }

}

export default Header;