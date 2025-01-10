import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './CreateRoom.module.css'

const CreateRoom = () => {
	const [hotels, setHotels] = useState([])
	const [formData, setFormData] = useState({
		hotel_id: '',
		room_type: '',
		price: '',
		amenities: '',
		capacity: '',
	})
	const [error, setError] = useState('')
	const navigate = useNavigate()

	useEffect(() => {
		fetchHotels()
	}, [])

	const fetchHotels = async () => {
		try {
			const response = await fetch(
				process.env.REACT_APP_API_URL + '/owners/hotels',
				{
					headers: {
						Authorization: `Bearer ${Cookies.get('token')}`,
					},
				}
			)
			if (response.ok) {
				const data = await response.json()
				setHotels(data)
			}
		} catch (err) {
			setError('Ошибка при загрузке отелей')
		}
	}

	const handleSubmit = async e => {
		e.preventDefault()

    const submitData = {
			...formData,
			hotel_id: parseInt(formData.hotel_id),
			price: parseFloat(formData.price),
			capacity: parseInt(formData.capacity),
		}
		try {
			const response = await fetch(
				process.env.REACT_APP_API_URL + '/owners/rooms',
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						Authorization: `Bearer ${Cookies.get('token')}`,
					},
					body: JSON.stringify(submitData),
				}
			)
			if (response.ok) {
				navigate('/owner/rooms')
			} else {
				const data = await response.json()
				setError(data.error)
			}
		} catch (err) {
			setError('Ошибка при создании номера')
		}
	}

	const handleChange = e => {
		const { name, value } = e.target
		setFormData(prev => ({
			...prev,
			[name]: value,
		}))
	}

	return (
		<div className={styles.container}>
			<button
				onClick={() => navigate('/owner/rooms')}
				className={styles.backButton}
			>
				Назад
			</button>
			<h1>Создание номера</h1>
			{error && <div className={styles.error}>{error}</div>}
			<form onSubmit={handleSubmit} className={styles.form}>
				<select
					name='hotel_id'
					value={formData.hotel_id}
					onChange={handleChange}
					required
				>
					<option value=''>Выберите отель</option>
					{hotels.map(hotel => (
						<option key={hotel.ID} value={hotel.ID}>
							{hotel.Name}
						</option>
					))}
				</select>
				<input
					type='text'
					name='room_type'
					placeholder='Тип номера'
					value={formData.room_type}
					onChange={handleChange}
					required
				/>
				<input
					type='number'
					name='price'
					placeholder='Цена за ночь'
					value={formData.price}
					onChange={handleChange}
					required
				/>
				<input
					type='text'
					name='amenities'
					placeholder='Удобства'
					value={formData.amenities}
					onChange={handleChange}
				/>
				<input
					type='number'
					name='capacity'
					placeholder='Вместимость'
					value={formData.capacity}
					onChange={handleChange}
					required
				/>
				<button type='submit'>Создать номер</button>
			</form>
		</div>
	)
}

export default CreateRoom
