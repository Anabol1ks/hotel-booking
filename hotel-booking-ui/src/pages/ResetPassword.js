import React, { useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './ResetPassword.module.css'

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
				const authToken = Cookies.get('token')
				setTimeout(() => {
					if (authToken) {
						navigate('/')
					} else {
						navigate('/auth/login')
					}
				}, 2000)
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
		<div className={styles.container}>
			<h2 className={styles.title}>Сброс пароля</h2>
			<p className={styles.description}>
				Введите новый пароль для вашей учётной записи.
			</p>
			{error && <p className={styles.error}>{error}</p>}
			{success && <p className={styles.success}>{success}</p>}
			<form onSubmit={handleSubmit}>
				<div className={styles.formGroup}>
					<label className={styles.label}>Новый пароль</label>
					<input
						className={styles.input}
						type='password'
						value={password}
						onChange={e => setPassword(e.target.value)}
						required
					/>
				</div>
				<div className={styles.formGroup}>
					<label className={styles.label}>Подтвердите пароль</label>
					<input
						className={styles.input}
						type='password'
						value={confirmPassword}
						onChange={e => setConfirmPassword(e.target.value)}
						required
					/>
				</div>
				<button className={styles.button} type='submit' disabled={loading}>
					{loading ? 'Изменение...' : 'Сменить пароль'}
				</button>
			</form>
		</div>
	)
}

export default ResetPassword
