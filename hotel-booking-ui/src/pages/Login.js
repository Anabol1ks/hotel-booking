import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'

const Login = () => {
	const navigate = useNavigate()
	const [formData, setFormData] = useState({
		email: '',
		password: '',
	})
	const [error, setError] = useState('')

	const handleChange = e => {
		const { name, value } = e.target
		setFormData(prev => ({ ...prev, [name]: value }))
	}

	const handleSubmit = async e => {
		e.preventDefault()
		setError('')

		try {
			const response = await fetch('http://localhost:8080/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(formData),
			})

			const data = await response.json()

			if (response.ok) {
				// Сохраняем токен в локальное хранилище
				localStorage.setItem('token', data.token)
				localStorage.setItem('role', data.role)
				navigate('/') // Перенаправление на главную страницу
			} else {
				setError(data.error || 'Ошибка авторизации')
			}
		} catch (err) {
			setError('Что-то пошло не так. Попробуйте позже.')
		}
	}

	return (
		<div>
			<h2>Вход</h2>
			{error && <p style={{ color: 'red' }}>{error}</p>}
			<form onSubmit={handleSubmit}>
				<div>
					<label>Email</label>
					<input
						type='email'
						name='email'
						value={formData.email}
						onChange={handleChange}
						required
					/>
				</div>
				<div>
					<label>Пароль</label>
					<input
						type='password'
						name='password'
						value={formData.password}
						onChange={handleChange}
						required
					/>
				</div>
				<button type='submit'>Войти</button>
			</form>
			<div style={{ marginTop: '10px' }}>
				<button
					type='button'
					onClick={() => navigate('/auth/forgot-password')}
					style={{
						background: 'none',
						border: 'none',
						color: 'blue',
						textDecoration: 'underline',
						cursor: 'pointer',
					}}
				>
					Забыли пароль?
				</button>
			</div>
		</div>
	)
}

export default Login
