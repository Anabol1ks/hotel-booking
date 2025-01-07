import React, { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './HotelsAndRooms.module.css'

const HotelsAndRooms = () => {
	const [hotels, setHotels] = useState([])
	const [rooms, setRooms] = useState([])
	const [error, setError] = useState('')
	const [isAuthenticated, setIsAuthenticated] = useState(false)
	const [favoriteRooms, setFavoriteRooms] = useState(new Set())
	const [filters, setFilters] = useState({
		minPrice: '',
		maxPrice: '',
		capacity: ''
	})
	const navigate = useNavigate()

	useEffect(() => {
		const token = Cookies.get('token')
		setIsAuthenticated(!!token)
	}, [])

	const fetchHotels = async () => {
		try {
			const response = await fetch('http://localhost:8080/hotels')
			if (response.ok) {
				const data = await response.json()
				setHotels(data)
			}
		} catch (err) {
			setError('Ошибка при загрузке отелей')
		}
	}

	const fetchRooms = useCallback(async () => {
		try {
			const queryParams = new URLSearchParams()
			if (filters.minPrice) queryParams.append('min_price', filters.minPrice)
			if (filters.maxPrice) queryParams.append('max_price', filters.maxPrice)
			if (filters.capacity) queryParams.append('capacity', filters.capacity)

			const response = await fetch(`http://localhost:8080/rooms?${queryParams}`)
			if (response.ok) {
				const data = await response.json()
				setRooms(data)
			}
		} catch (err) {
			setError('Ошибка при загрузке номеров')
		}
	}, [filters])

	const fetchFavorites = async () => {
		if (!isAuthenticated) return
		
		try {
			const response = await fetch('http://localhost:8080/favorites', {
				headers: {
					Authorization: `Bearer ${Cookies.get('token')}`
				}
			})
			if (response.ok) {
				const data = await response.json()
				setFavoriteRooms(new Set(data.map(room => room.ID)))
			}
		} catch (err) {
			console.error('Ошибка при загрузке избранного:', err)
		}
	}

	const handleAddToFavorites = async (roomId) => {
		if (!isAuthenticated) {
			alert('Пожалуйста, авторизуйтесь чтобы добавить номер в избранное')
			navigate('/login')
			return
		}

		try {
			const response = await fetch(`http://localhost:8080/favorites/${roomId}`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${Cookies.get('token')}`,
				},
			})

			if (response.ok) {
				setFavoriteRooms(prev => new Set([...prev, roomId]))
				alert('Номер добавлен в избранное')
			}
		} catch (err) {
			console.error('Ошибка при добавлении в избранное:', err)
		}
	}

	const removeFromFavorites = async (roomId) => {
		try {
			const response = await fetch(`http://localhost:8080/favorites/${roomId}`, {
				method: 'DELETE',
				headers: {
					Authorization: `Bearer ${Cookies.get('token')}`
				}
			})
			if (response.ok) {
				setFavoriteRooms(prev => {
					const newSet = new Set(prev)
					newSet.delete(roomId)
					return newSet
				})
				alert('Номер удален из избранного')
			}
		} catch (err) {
			console.error('Ошибка при удалении из избранного:', err)
		}
	}

	useEffect(() => {
		fetchHotels()
	}, [])

	useEffect(() => {
		fetchRooms()
	}, [fetchRooms])

	useEffect(() => {
		fetchFavorites()
	}, [isAuthenticated])

	const handleFilterChange = (e) => {
		const { name, value } = e.target
		setFilters(prev => ({
			...prev,
			[name]: value
		}))
	}

	return (
		<div className={styles.container}>
			<div>
				<button className={styles.backButton} onClick={() => navigate('/')}>
					Вернуться на главную
				</button>
			</div>
			<h1 className={styles.title}>Отели и номера</h1>
			<div className={styles.filters}>
				<h3>Фильтры</h3>
				<div className={styles.filterGroup}>
					<input
						type='number'
						name='minPrice'
						placeholder='Мин. цена'
						value={filters.minPrice}
						onChange={handleFilterChange}
						className={styles.filterInput}
					/>
					<input
						type='number'
						name='maxPrice'
						placeholder='Макс. цена'
						value={filters.maxPrice}
						onChange={handleFilterChange}
						className={styles.filterInput}
					/>
					<input
						type='number'
						name='capacity'
						placeholder='Количество гостей'
						value={filters.capacity}
						onChange={handleFilterChange}
						className={styles.filterInput}
					/>
				</div>
			</div>

			{error && <div className={styles.error}>{error}</div>}

			<div className={styles.hotelsSection}>
				<h2>Отели</h2>
				<div className={styles.hotelsList}>
					{hotels.map(hotel => (
						<div key={hotel.ID} className={styles.hotelCard}>
							<h3>{hotel.Name}</h3>
							<p>{hotel.Address}</p>
							<p className={styles.description}>{hotel.Description}</p>
						</div>
					))}
				</div>
			</div>

			<div className={styles.roomsSection}>
				<h2>Доступные номера</h2>
				<div className={styles.roomsList}>
					{rooms.map(room => (
						<div key={room.ID} className={styles.roomCard}>
							<h3>Тип номера: {room.RoomType}</h3>
							<p>Цена: {room.Price} руб/ночь</p>
							<p>Вместимость: {room.Capacity} чел.</p>
							<p>Удобства: {room.Amenities}</p>
							{favoriteRooms.has(room.ID) ? (
								<button 
									onClick={() => removeFromFavorites(room.ID)}
									className={`${styles.favoriteButton} ${styles.removeButton}`}
								>
									Убрать из избранного
								</button>
							) : (
								<button 
									onClick={() => handleAddToFavorites(room.ID)}
									className={styles.favoriteButton}
								>
									Добавить в избранное
								</button>
							)}
						</div>
					))}
				</div>
			</div>
		</div>
	)
}

export default HotelsAndRooms