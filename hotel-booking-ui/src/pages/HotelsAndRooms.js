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
	const [bookingDates, setBookingDates] = useState({})
	const [filters, setFilters] = useState({
		minPrice: '',
		maxPrice: '',
		capacity: '',
		startDate: '',
		endDate: '',
		hotelId: ''
	})
	const navigate = useNavigate()

	useEffect(() => {
		const token = Cookies.get('token')
		setIsAuthenticated(!!token)
	}, [])

	const fetchHotels = async () => {
		try {
			const response = await fetch(process.env.REACT_APP_API_URL + '/hotels')
			if (response.ok) {
				const data = await response.json()
				setHotels(data)
			}
		} catch (err) {
			setError('Ошибка при загрузке отелей')
		}
	}

	const [currentImageIndexes, setCurrentImageIndexes] = useState({})

	useEffect(() => {
		const initialIndexes = {}
		rooms.forEach(room => {
			initialIndexes[room.ID] = 0
		})
		setCurrentImageIndexes(initialIndexes)
	}, [rooms])

	const nextImage = (roomId, totalImages) => {
		setCurrentImageIndexes(prev => ({
			...prev,
			[roomId]: prev[roomId] === undefined ? 0 : (prev[roomId] + 1) % totalImages,
		}))
	}

	const prevImage = (roomId, totalImages) => {
		setCurrentImageIndexes(prev => ({
			...prev,
			[roomId]:
				prev[roomId] === undefined
					? 0
					: (prev[roomId] - 1 + totalImages) % totalImages,
		}))
	}

	const fetchRooms = useCallback(async () => {
		try {
			const queryParams = new URLSearchParams()
			if (filters.minPrice) queryParams.append('min_price', filters.minPrice)
			if (filters.maxPrice) queryParams.append('max_price', filters.maxPrice)
			if (filters.capacity) queryParams.append('capacity', filters.capacity)
			if (filters.startDate) queryParams.append('start_date', filters.startDate)
			if (filters.endDate) queryParams.append('end_date', filters.endDate)
			if (filters.hotelId) queryParams.append('hotel_id', filters.hotelId)

			const response = await fetch(
				process.env.REACT_APP_API_URL + `/rooms?${queryParams}`
			)
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
			const response = await fetch(
				process.env.REACT_APP_API_URL + '/favorites',
				{
					headers: {
						Authorization: `Bearer ${Cookies.get('token')}`,
					},
				}
			)
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
			navigate('auth/login')
			return
		}

		try {
			const response = await fetch(process.env.REACT_APP_API_URL + `/favorites/${roomId}`, {
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
			const response = await fetch(process.env.REACT_APP_API_URL + `/favorites/${roomId}`, {
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

	const handleFilterChange = e => {
		const { name, value } = e.target

		if (name === 'startDate') {
			const today = new Date().toISOString().split('T')[0]
			if (value < today) {
				alert('Дата не может быть в прошлом')
				return
			}
			if (filters.endDate && value > filters.endDate) {
				alert('Дата начала не может быть позже даты окончания')
				return
			}
		}

		if (name === 'endDate') {
			if (filters.startDate && value < filters.startDate) {
				alert('Дата окончания не может быть раньше даты начала')
				return
			}
		}

		setFilters(prev => ({
			...prev,
			[name]: value,
		}))
	}


	const handleBookRoom = async roomId => {
		if (!isAuthenticated) {
			alert('Пожалуйста, авторизуйтесь чтобы забронировать номер')
			navigate('/login')
			return
		}

		const startDate = bookingDates[roomId]?.startDate
		const endDate = bookingDates[roomId]?.endDate

		if (!startDate || !endDate) {
			alert('Пожалуйста, выберите даты бронирования')
			return
		}

		try {
			const response = await fetch(
				process.env.REACT_APP_API_URL + '/bookings',
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						Authorization: `Bearer ${Cookies.get('token')}`,
					},
					body: JSON.stringify({
						room_id: roomId,
						start_date: new Date(startDate).toISOString(),
						end_date: new Date(endDate).toISOString(),
					}),
				}
			)

			if (response.ok) {
				alert('Номер успешно забронирован!')
				setBookingDates(prev => ({
					...prev,
					[roomId]: { startDate: null, endDate: null },
				}))
			} else {
				const errorData = await response.json()
				alert(errorData.error || 'Ошибка при бронировании')
			}
		} catch (err) {
			console.error('Ошибка при бронировании:', err)
			alert('Не удалось забронировать номер')
		}
	}

	const handleDateChange = (roomId, dateType, value) => {
		const today = new Date().toISOString().split('T')[0]

		// Validate start date
		if (dateType === 'startDate') {
			if (value < today) {
				alert('Дата не может быть в прошлом')
				return
			}

			// Check if end date exists and start date is after end date
			const currentEndDate = bookingDates[roomId]?.endDate
			if (currentEndDate && value > currentEndDate) {
				alert('Дата начала не может быть позже даты окончания')
				return
			}
		}

		// Validate end date
		if (dateType === 'endDate') {
			const currentStartDate = bookingDates[roomId]?.startDate
			if (currentStartDate && value < currentStartDate) {
				alert('Дата окончания не может быть раньше даты начала')
				return
			}
		}

		// If all validations pass, update booking dates
		setBookingDates(prev => ({
			...prev,
			[roomId]: {
				...prev[roomId],
				[dateType]: value,
			},
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
					<select
						name='hotelId'
						value={filters.hotelId}
						onChange={handleFilterChange}
						className={styles.filterInput}
					>
						<option value=''>Выберите отель</option>
						{hotels.map(hotel => (
							<option key={hotel.ID} value={hotel.ID}>
								{hotel.Name}
							</option>
						))}
					</select>
					<input
						type='date'
						name='startDate'
						value={filters.startDate}
						onChange={handleFilterChange}
						className={styles.filterInput}
						min={new Date().toISOString().split('T')[0]}
					/>
					<input
						type='date'
						name='endDate'
						value={filters.endDate}
						onChange={handleFilterChange}
						className={styles.filterInput}
						min={filters.startDate || new Date().toISOString().split('T')[0]}
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
							<div className={styles.roomImages}>
								{room.Images && room.Images.length > 0 ? (
									<div className={styles.imageCarousel}>
										{room.Images.length > 1 && (
											<button
												className={`${styles.arrowButton} ${styles.leftArrow}`}
												onClick={() => prevImage(room.ID, room.Images.length)}
											>
												←
											</button>
										)}
										<img
											src={
												room.Images[currentImageIndexes[room.ID] || 0]?.ImageURL
											}
											alt={`Room ${room.RoomType}`}
											className={styles.roomImage}
										/>
										{room.Images.length > 1 && (
											<button
												className={`${styles.arrowButton} ${styles.rightArrow}`}
												onClick={() => nextImage(room.ID, room.Images.length)}
											>
												→
											</button>
										)}
									</div>
								) : (
									<div className={styles.noImage}>
										Нет доступных изображений
									</div>
								)}
							</div>
							<h3>Тип номера: {room.RoomType}</h3>
							<p>Цена: {room.Price} руб/ночь</p>
							<p>Вместимость: {room.Capacity} чел.</p>
							<p>Удобства: {room.Amenities}</p>

							{isAuthenticated && (
								<div className={styles.bookingSection}>
									<div className={styles.bookingDates}>
										<label>Дата заезда:</label>
										<input
											type='date'
											value={bookingDates[room.ID]?.startDate || ''}
											onChange={e =>
												handleDateChange(room.ID, 'startDate', e.target.value)
											}
										/>
										<label>Дата выезда:</label>
										<input
											type='date'
											value={bookingDates[room.ID]?.endDate || ''}
											onChange={e =>
												handleDateChange(room.ID, 'endDate', e.target.value)
											}
										/>
									</div>
									<button
										onClick={() => handleBookRoom(room.ID)}
										className={styles.bookButton}
									>
										Забронировать
									</button>
								</div>
							)}

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

