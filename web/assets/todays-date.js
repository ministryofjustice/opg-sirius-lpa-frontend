const todaysDate = () => {
    let selectTodayLink = document.querySelector('[data-module="select-todays-date"]');
    let datePicker = document.querySelector('.date-picker');

    selectTodayLink.addEventListener(
        "click",
        function () {
            datePicker.value = new Date().toJSON().slice(0, 10)
        }, false
    );
};
export default todaysDate;
