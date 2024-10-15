document.addEventListener('DOMContentLoaded', function() {
    // 日期选择器初始化
    flatpickr(".datepicker", {
        dateFormat: "Y-m-d",
        locale: "zh"
    });

    // 工时提交表单验证
    const timesheetForm = document.getElementById('timesheet-form');
    if (timesheetForm) {
        timesheetForm.addEventListener('submit', function(e) {
            const hours = document.getElementById('hours').value;
            const date = document.getElementById('date').value;
            const project = document.getElementById('project').value;

            if (!hours || !date || !project) {
                e.preventDefault();
                alert('请填写所有必填字段');
            }
        });
    }

    // 项目选择器
    const projectSelect = document.getElementById('project');
    if (projectSelect) {
        new Choices(projectSelect, {
            searchEnabled: true,
            itemSelectText: ''
        });
    }
});
