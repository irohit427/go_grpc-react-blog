import React from 'react'
import { Form, Input, Button } from 'antd';
import './Signup.scss';
import { SignupRequest, AuthUserRequest } from "../../proto/services_grpc_web_pb";
import { authClient } from '../../App';

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
    let { username, email, password } = values;
    let request = new SignupRequest()
    request.setUsername(username)
    request.setEmail(email)
    request.setPassword(password)

    authClient.signup(request, {}, (err, res) => {
      if(err) return alert(err)
      localStorage.setItem('token', res.getToken())
      request = new AuthUserRequest()
      request.setToken(res.getToken())
      authClient.authUser(request, {}, (err, res) => {
          if(err) return alert(err)
          const user = { id: res.getId(), username: res.getUsername(), email: res.getEmail() }
          localStorage.setItem("user", JSON.stringify(user))
      })
    })
  };

  const onFinishFailed = (errorInfo) => {
    console.log('Failed:', errorInfo);
  };

  return (
    <div className="signupContainer">
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
          label="Email"
          name="email"
          rules={[
            {
              required: true,
              message: 'Please input your email!',
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

        <Form.Item {...tailLayout}>
          <Button type="primary" htmlType="submit">
            Sign Up
          </Button>
        </Form.Item>
      </Form>
    </div>
  )
}
