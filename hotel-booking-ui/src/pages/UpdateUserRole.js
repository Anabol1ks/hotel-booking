import React, { useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './UpdateUserRole.module.css'

const UpdateUserRole = () => {
	const location = useLocation()
	const navigate = useNavigate()
	const user = location.state?.user
	const [role, setRole] = useState(user?.role || '')
	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')
	const [loading, setLoading] = useState(false)
	const roles = ['client', 'owner', 'admin']

	const handleSubmit = async e => {
		e.preventDefault()
		setLoading(true)
		setError('')
		setSuccess('')

		try {
			const token = Cookies.get('token')
			const response = await fetch(
				`http://localhost:8080/admin/users/${user.ID}/role`,
				{
					method: 'PUT',
					headers: {
						'Content-Type': 'application/json',
						Authorization: `Bearer ${token}`,
					},
					body: JSON.stringify({ role }),
				}
			)

			if (response.ok) {
				setSuccess('Роль успешно обновлена')
			} else {
				const data = await response.json()
				setError(data.error || 'Ошибка при обновлении роли')
			}
		} catch (err) {
			setError('Произошла ошибка. Попробуйте еще раз.')
		} finally {
			setLoading(false)
		}
	}

	if (!user) {
		return <p>Данные о пользователе не переданы</p>
	}

	return (
		<div className={styles.container}>
			<h1 className={styles.title}>Смена роли пользователя</h1>
			<div className={styles.userInfo}>
				<p>Редактирование пользователя: {user.Name}</p>
				<p>Email: {user.Email}</p>
			</div>

			{error && <p className={styles.error}>{error}</p>}
			{success && <p className={styles.success}>{success}</p>}

			<form onSubmit={handleSubmit}>
				<div className={styles.formGroup}>
					<label className={styles.label} htmlFor='role'>
						Выберите новую роль:
					</label>
					<select
						id='role'
						className={styles.select}
						value={role}
						onChange={e => setRole(e.target.value)}
						required
					>
						<option value='' disabled>
							Выберите роль
						</option>
						{roles.map(r => (
							<option key={r} value={r}>
								{r}
							</option>
						))}
					</select>
				</div>
				<button
					type='submit'
					className={`${styles.button} ${styles.submitButton}`}
					disabled={loading}
				>
					{loading ? 'Обновление...' : 'Обновить роль'}
				</button>
			</form>
			<button
				onClick={() => navigate('/admin/users')}
				className={`${styles.button} ${styles.backButton}`}
			>
				Вернуться назад
			</button>
		</div>
	)
}

export default UpdateUserRole
