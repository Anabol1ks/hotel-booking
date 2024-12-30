import React, { useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'

const ResetPassword = () => {
	const [password, setPassword] = useState('')
	const [confirmPassword, setConfirmPassword] = useState('')
	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')
	const [loading, setLoading] = useState(false)
	const navigate = useNavigate()
	const [searchParams] = useSearchParams()

	const token = searchParams.get('token')

	const handleSubmit = async e => {
		e.preventDefault()
		setError('')
		setSuccess('')
		setLoading(true)

		if (!token) {
			setError('Токен для сброса пароля не найден.')
			setLoading(false)
			return
		}

		if (password !== confirmPassword) {
			setError('Пароли не совпадают.')
			setLoading(false)
			return
		}

		try {
			const response = await fetch(
				'http://localhost:8080/auth/reset-password',
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({ token, password }),
				}
			)

			if (response.ok) {
				setSuccess('Пароль успешно изменён.')

				// Проверка наличия токена авторизации
				const authToken = localStorage.getItem('token')
				setTimeout(() => {
					if (authToken) {
						navigate('/') // Перенаправление на аккаунт
					} else {
						navigate('/auth/login') // Перенаправление на страницу входа
					}
				}, 2000) // Небольшая задержка для отображения сообщения успеха
			} else {
				const data = await response.json()
				setError(data.error || 'Ошибка при сбросе пароля.')
			}
		} catch (err) {
			setError('Произошла ошибка. Попробуйте ещё раз позже.')
		} finally {
			setLoading(false)
		}
	}

	return (
		<div>
			<h2>Сброс пароля</h2>
			<p>Введите новый пароль для вашей учётной записи.</p>
			{error && <p style={{ color: 'red' }}>{error}</p>}
			{success && <p style={{ color: 'green' }}>{success}</p>}
			<form onSubmit={handleSubmit}>
				<div>
					<label>Новый пароль</label>
					<input
						type='password'
						value={password}
						onChange={e => setPassword(e.target.value)}
						required
					/>
				</div>
				<div>
					<label>Подтвердите пароль</label>
					<input
						type='password'
						value={confirmPassword}
						onChange={e => setConfirmPassword(e.target.value)}
						required
					/>
				</div>
				<button type='submit' disabled={loading}>
					{loading ? 'Изменение...' : 'Сменить пароль'}
				</button>
			</form>
		</div>
	)
}

export default ResetPassword
