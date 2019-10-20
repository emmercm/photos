import React, {Component} from 'react';
import Albums from "./components/Albums";
import './App.css'

import { Layout, Menu } from 'antd';

const { Header, Content, Footer, Sider } = Layout;

class App extends Component {
    state = {
        albums: [
            {
                id: 1,
                path: '/',
                title: 'ABC'
            },
            {
                id: 2,
                path: '/',
                title: 'DEF'
            },
            {
                id: 3,
                path: '/',
                title: 'GHI'
            }
        ]
    };

    render() {
        return (
            <Layout>
                <Header></Header>
                <Content>
                    <Layout>
                        <Sider>
                            <Menu>
                                <Menu.Item key="1">ok</Menu.Item>
                                <Menu.Item key="2">ok</Menu.Item>
                                <Menu.Item key="3">ok</Menu.Item>
                                <Menu.Item key="4">ok</Menu.Item>
                                <Menu.Item key="5">ok</Menu.Item>
                                <Menu.Item key="6">ok</Menu.Item>
                                <Menu.Item key="7">ok</Menu.Item>
                                <Menu.Item key="8">ok</Menu.Item>
                                <Menu.Item key="9">ok</Menu.Item>
                                <Menu.Item key="10">ok</Menu.Item>
                                <Menu.Item key="11">ok</Menu.Item>
                                <Menu.Item key="12">ok</Menu.Item>
                                <Menu.Item key="13">ok</Menu.Item>
                                <Menu.Item key="14">ok</Menu.Item>
                                <Menu.Item key="15">ok</Menu.Item>
                                <Menu.Item key="16">ok</Menu.Item>
                                <Menu.Item key="17">ok</Menu.Item>
                            </Menu>
                        </Sider>
                        <Content>
                            <Albums albums={this.state.albums}/>
                        </Content>
                    </Layout>
                </Content>
                <Footer></Footer>
            </Layout>
        );
    }
}


export default App;
