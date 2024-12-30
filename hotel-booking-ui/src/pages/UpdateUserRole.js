import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'

const UpdateUserRole = () => {
	const { id } = useParams() // Получение ID пользователя из URL
	const [role, setRole] = useState('')
	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')
	const [loading, setLoading] = useState(false)
	const navigate = useNavigate()

	// Список доступных ролей
	const roles = ['client', 'owner', 'admin']

	const handleSubmit = async e => {
		e.preventDefault()

		setLoading(true)
		setError('')
		setSuccess('')

		try {
			const token = localStorage.getItem('token')

			const response = await fetch(
				`http://localhost:8080/admin/users/${id}/role`,
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

	useEffect(() => {
		// Дополнительно можно получить текущую роль пользователя, чтобы отобразить её
		const fetchUserRole = async () => {
			try {
				const token = localStorage.getItem('token')
				const response = await fetch(
					`http://localhost:8080/admin/users/${id}`,
					{
						headers: {
							Authorization: `Bearer ${token}`,
						},
					}
				)

				if (response.ok) {
					const data = await response.json()
					setRole(data.role) // Устанавливаем текущую роль пользователя
				} else {
					setError('Не удалось загрузить данные пользователя')
				}
			} catch (err) {
				setError('Произошла ошибка при загрузке данных')
			}
		}

		fetchUserRole()
	}, [id])

	return (
		<div>
			<h1>Смена роли пользователя</h1>
			{id && <p>Редактирование пользователя с ID: {id}</p>}

			{error && <p style={{ color: 'red' }}>{error}</p>}
			{success && <p style={{ color: 'green' }}>{success}</p>}

			<form onSubmit={handleSubmit}>
				<div>
					<label htmlFor='role'>Выберите новую роль:</label>
					<select
						id='role'
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
				<button type='submit' disabled={loading}>
					{loading ? 'Обновление...' : 'Обновить роль'}
				</button>
			</form>
			<button onClick={() => navigate('/admin/users')}>Вернуться назад</button>
		</div>
	)
}

export default UpdateUserRole
