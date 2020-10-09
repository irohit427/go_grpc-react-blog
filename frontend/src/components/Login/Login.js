import React from 'react'
import { Form, Input, Button, Checkbox } from 'antd';
import './Login.scss';
import { authClient } from '../../App';
import { LoginRequest, AuthUserRequest } from "../../proto/services_grpc_web_pb";



const layout = {
  labelCol: {
    span: 8,
  },
  wrapperCol: {
    span: 16,
  },
};
const tailLayout = {
  wrapperCol: {
    offset: 1,
    span: 16,
  },
};

export default function Login() {
  const onFinish = (values) => {
    let { username, password } = values;
    let req = new LoginRequest()
    req.setUsername(username)
    req.setPassword(password)
    authClient.login(req, {}, (err, res) => {
      if(err) return alert(err.message)
      localStorage.setItem('token', res.getToken())
      req = new AuthUserRequest()
      req.setToken(res.getToken())
      authClient.authUser(req, {}, (err, res) => {
          if(err) return alert(err.message)
          const user = { id: res.getId(), username: res.getUsername(), email: res.getEmail() }
          localStorage.setItem('user', JSON.stringify(user))
      })
      console.log(res.getToken());
  })
  };

  const onFinishFailed = (errorInfo) => {
    console.log('Failed:', errorInfo);
  };

  return (
    <div className="loginContainer">
      <Form
        {...layout}
        name="basic"
        initialValues={{
          remember: true,
        }}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
      >
        <Form.Item
          label="Username"
          name="username"
          rules={[
            {
              required: true,
              message: 'Please input your username!',
            },
          ]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Password"
          name="password"
          rules={[
            {
              required: true,
              message: 'Please input your password!',
            },
          ]}
        >
          <Input.Password />
        </Form.Item>

        <Form.Item {...tailLayout} name="remember" valuePropName="checked">
          <Checkbox>Remember me</Checkbox>
        </Form.Item>

        <Form.Item {...tailLayout}>
          <Button type="primary" htmlType="submit">
            Submit
          </Button>
        </Form.Item>
      </Form>
    </div>
  )
}
