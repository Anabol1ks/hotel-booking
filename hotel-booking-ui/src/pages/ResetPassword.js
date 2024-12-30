import React, { useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'

const ResetPassword = () => {
	const location = useLocation()
	const history = useNavigate()

	const [password, setPassword] = useState('')
	const [error, setError] = useState('')

	// Функция для извлечения параметров URL
	const getQueryParam = name => {
		const urlParams = new URLSearchParams(location.search)
		return urlParams.get(name)
	}

	const token = getQueryParam('token')

	const handleSubmit = async e => {
		e.preventDefault()

		if (!password) {
			setError('Password is required')
			return
		}

		const body = { token, password }

		try {
			const response = await fetch(
				'http://localhost:8080/auth/reset-password',
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify(body),
				}
			)

			const data = await response.json() // Распарсим ответ

			if (response.ok) {
				// Успешный сброс пароля
				history('/login')
			} else {
				// Установка ошибки с сервера
				setError(data.error || 'Failed to reset password')
			}
		} catch (err) {
			setError('Something went wrong')
		}
	}

	return (
		<div>
			<h2>Reset Password</h2>
			{error && <p style={{ color: 'red' }}>{error}</p>}
			<form onSubmit={handleSubmit}>
				<div>
					<label>Password</label>
					<input
						type='password'
						value={password}
						onChange={e => setPassword(e.target.value)}
						required
					/>
				</div>
				<button type='submit'>Reset Password</button>
			</form>
		</div>
	)
}

console.log('ResetPassword component loaded')
export default ResetPassword
