import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'

const Register = () => {
	const history = useNavigate()
	const [formData, setFormData] = useState({
		name: '',
		email: '',
		password: '',
		phone: '',
	})
	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')

	const handleChange = e => {
		const { name, value } = e.target
		setFormData(prev => ({ ...prev, [name]: value }))
	}

	const handleSubmit = async e => {
		e.preventDefault()
		setError('')
		setSuccess('')

		try {
			const response = await fetch(
				process.env.REACT_APP_API_URL + '/auth/register',
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify(formData),
				}
			)

			const data = await response.json()

			if (response.status === 201) {
				setSuccess('Регистрация успешна!')
				history('auth/login') // Перенаправление на страницу входа
			} else {
				setError(data.error || 'Ошибка регистрации')
			}
		} catch (err) {
			setError('Что-то пошло не так. Попробуйте позже.')
		}
	}

	return (
		<div>
			<h2>Регистрация</h2>
			{error && <p style={{ color: 'red' }}>{error}</p>}
			{success && <p style={{ color: 'green' }}>{success}</p>}
			<form onSubmit={handleSubmit}>
				<div>
					<label>Имя</label>
					<input
						type='text'
						name='name'
						value={formData.name}
						onChange={handleChange}
						required
					/>
				</div>
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
				<div>
					<label>Телефон</label>
					<input
						type='text'
						name='phone'
						value={formData.phone}
						onChange={handleChange}
						required
					/>
				</div>
				<button type='submit'>Зарегистрироваться</button>
			</form>
		</div>
	)
}

export default Register
