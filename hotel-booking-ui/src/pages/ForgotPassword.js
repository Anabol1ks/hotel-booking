import React, { useState } from 'react'
import styles from './ForgotPassword.module.css'

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
				setEmail('')
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
		<div className={styles.container}>
			<h2 className={styles.title}>Восстановление пароля</h2>
			<p className={styles.description}>
				Введите свой email, и мы отправим вам письмо с инструкцией по сбросу
				пароля.
			</p>
			{error && <p className={styles.error}>{error}</p>}
			{success && <p className={styles.success}>{success}</p>}
			<form onSubmit={handleSubmit}>
				<div className={styles.formGroup}>
					<label className={styles.label}>Email</label>
					<input
						className={styles.input}
						type="email"
						value={email}
						onChange={e => setEmail(e.target.value)}
						required
					/>
				</div>
				<button 
					className={styles.button} 
					type="submit" 
					disabled={loading}
				>
					{loading ? 'Отправка...' : 'Отправить'}
				</button>
			</form>
		</div>
	)
}

export default ForgotPassword