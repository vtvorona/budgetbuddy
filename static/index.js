// Открыть модальное окно добавления
function toggleAddModal() {
    const modal = document.getElementById("addmodal");
    modal.classList.toggle("hidden");
    modal.classList.toggle("flex");
}

// Закрыть модальное окно добавления
function closeToggleAddModal() {
    const modal = document.getElementById("addmodal");
    modal.classList.add("hidden");
    modal.classList.remove("flex");
}

// Закрыть модальное окно дня
function closeDayModal() {
    const modal = document.getElementById("day-modal");
    modal.classList.add("hidden");
    modal.classList.remove("flex");
}

// Открыть модальное окно редактирования
function openEditModal(expense) {
    const modal = document.getElementById("edit-modal");

    // Заполнение полей формы
    document.getElementById("edit-expense-id").value = expense.id || "";
    document.getElementById("edit-title").value = expense.title || "";
    document.getElementById("edit-category").value = expense.categoryId || ""; // Установить текущую категорию
    document.getElementById("edit-price").value = expense.price || "";
    document.getElementById("edit-amount").value = expense.amount || "";

    modal.classList.remove("hidden");
    modal.classList.add("flex");
}

// Закрыть модальное окно редактирования
function closeEditModal() {
    const modal = document.getElementById("edit-modal");
    modal.classList.add("hidden");
    modal.classList.remove("flex");
}

// Добавить кнопку редактирования к расходу
function createEditButton(expense) {
    const editButton = document.createElement("button");
    editButton.textContent = "Редактировать";
    editButton.classList.add("px-4", "py-2", "bg-blue-500", "text-white", "rounded-lg", "hover:bg-blue-600");
    editButton.dataset.id = expense.dataset.id;
    editButton.dataset.title = expense.dataset.title;
    editButton.dataset.categoryId = expense.dataset.categoryId;
    editButton.dataset.price = expense.dataset.price;
    editButton.dataset.amount = expense.dataset.amount;

    // Добавить обработчик события
    editButton.addEventListener("click", () => {
        openEditModal({
            id: editButton.dataset.id,
            title: editButton.dataset.title,
            categoryId: editButton.dataset.categoryId,
            price: editButton.dataset.price,
            amount: editButton.dataset.amount,
        });
    });

    return editButton;
}

// Открыть модальное окно дня
function openDayModal(date) {
    const modal = document.getElementById(`day-modal-${date}`);
    if (modal) {
        modal.classList.remove("hidden");
        modal.classList.add("flex");
    } else {
        console.error(`Modal for date ${date} not found`);
    }
}

function closeDayModal(date) {
    const modal = document.getElementById(`day-modal-${date}`);
    if (modal) {
        modal.classList.add("hidden");
        modal.classList.remove("flex");
    } else {
        console.error(`Modal for date ${date} not found`);
    }
}
