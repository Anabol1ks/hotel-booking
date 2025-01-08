import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './CreateOfflineBooking.module.css'

const CreateOfflineBooking = () => {
	const navigate = useNavigate()
	const [formData, setFormData] = useState({
		room_id: '',
		start_date: '',
		end_date: '',
		phone_number: '',
		name: '',
	})
	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')
	const [loading, setLoading] = useState(false)

	const handleChange = e => {
		const { name, value } = e.target
		setFormData(prev => ({
			...prev,
			[name]: value,
		}))
	}

	const handleSubmit = async e => {
		e.preventDefault()
		setLoading(true)
		setError('')
		setSuccess('')

		try {
			const token = Cookies.get('token')
			const response = await fetch('http://localhost:8080/booking/offline', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`,
				},
				body: JSON.stringify({
					...formData,
					start_date: new Date(formData.start_date).toISOString(),
					end_date: new Date(formData.end_date).toISOString(),
					room_id: parseInt(formData.room_id),
				}),
			})

			const data = await response.json()

			if (response.ok) {
				setSuccess('Бронирование успешно создано')
				setTimeout(() => navigate('/bookings'), 2000)
			} else {
				setError(data.error || 'Ошибка при создании бронирования')
			}
		} catch (err) {
			setError('Произошла ошибка при создании бронирования')
		} finally {
			setLoading(false)
		}
	}

	return (
		<div className={styles.container}>
			<h1>Создание офлайн бронирования</h1>

			{error && <div className={styles.error}>{error}</div>}
			{success && <div className={styles.success}>{success}</div>}

			<form onSubmit={handleSubmit} className={styles.form}>
				<div className={styles.formGroup}>
					<label>Номер комнаты:</label>
					<input
						type='number'
						name='room_id'
						value={formData.room_id}
						onChange={handleChange}
						required
					/>
				</div>

				<div className={styles.formGroup}>
					<label>Имя клиента:</label>
					<input
						type='text'
						name='name'
						value={formData.name}
						onChange={handleChange}
						required
					/>
				</div>

				<div className={styles.formGroup}>
					<label>Номер телефона:</label>
					<input
						type='tel'
						name='phone_number'
						value={formData.phone_number}
						onChange={handleChange}
						required
					/>
				</div>

				<div className={styles.formGroup}>
					<label>Дата заезда:</label>
					<input
						type='datetime-local'
						name='start_date'
						value={formData.start_date}
						onChange={handleChange}
						required
					/>
				</div>

				<div className={styles.formGroup}>
					<label>Дата выезда:</label>
					<input
						type='datetime-local'
						name='end_date'
						value={formData.end_date}
						onChange={handleChange}
						required
					/>
				</div>

				<button
					type='submit'
					className={styles.submitButton}
					disabled={loading}
				>
					{loading ? 'Создание...' : 'Создать бронирование'}
				</button>
			</form>
		</div>
	)
}

export default CreateOfflineBooking
