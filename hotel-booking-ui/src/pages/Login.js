import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './Login.module.css'

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
				Cookies.set('token', data.token)
				Cookies.set('role', data.role)
				navigate('/')
			} else {
				setError(data.error || 'Ошибка авторизации')
			}
		} catch (err) {
			setError('Что-то пошло не так. Попробуйте позже.')
		}
	}

	return (
		<div className={styles.loginContainer}>
			<h2 className={styles.title}>Вход</h2>
			{error && <p className={styles.error}>{error}</p>}
			<form onSubmit={handleSubmit}>
				<div className={styles.formGroup}>
					<label className={styles.label}>Email</label>
					<input
						className={styles.input}
						type='email'
						name='email'
						value={formData.email}
						onChange={handleChange}
						required
					/>
				</div>
				<div className={styles.formGroup}>
					<label className={styles.label}>Пароль</label>
					<input
						className={styles.input}
						type='password'
						name='password'
						value={formData.password}
						onChange={handleChange}
						required
					/>
				</div>
				<button className={styles.submitButton} type='submit'>
					Войти
				</button>
			</form>
			<div className={styles.forgotPassword}>
				<button
					className={styles.forgotPasswordButton}
					type='button'
					onClick={() => navigate('/auth/forgot-password')}
				>
					Забыли пароль?
				</button>
			</div>
		</div>
	)
}

export default Login
