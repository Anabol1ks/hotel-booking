import React, { useState } from 'react'

const ForgotPassword = () => {
	const [email, setEmail] = useState('')
	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')
	const [loading, setLoading] = useState(false)

	const handleSubmit = async e => {
		e.preventDefault()
		setError('')
		setSuccess('')
		setLoading(true)

		try {
			const response = await fetch(
				'http://localhost:8080/auth/reset-password-request',
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({ email }),
				}
			)

			if (response.ok) {
				setSuccess('Письмо с инструкциями отправлено на ваш email.')
				setEmail('') // Сброс email
			} else {
				const data = await response.json()
				setError(data.error || 'Ошибка при отправке письма.')
			}
		} catch (err) {
			setError('Произошла ошибка. Попробуйте ещё раз позже.')
		} finally {
			setLoading(false)
		}
	}

	return (
		<div>
			<h2>Восстановление пароля</h2>
			<p>
				Введите свой email, и мы отправим вам письмо с инструкцией по сбросу
				пароля.
			</p>
			{error && <p style={{ color: 'red' }}>{error}</p>}
			{success && <p style={{ color: 'green' }}>{success}</p>}
			<form onSubmit={handleSubmit}>
				<div>
					<label>Email</label>
					<input
						type='email'
						value={email}
						onChange={e => setEmail(e.target.value)}
						required
					/>
				</div>
				<button type='submit' disabled={loading}>
					{loading ? 'Отправка...' : 'Отправить'}
				</button>
			</form>
		</div>
	)
}

export default ForgotPassword
