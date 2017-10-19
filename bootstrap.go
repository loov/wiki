package main

var Welcome = Page{
	URL:   "http://localhost:8000/welcome",
	Title: "Welcome",
	Story: []Item{
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Optio rerum nam architecto [[mollitia]] unde officia tempora reiciendis omnis quas, expedita, quia ad culpa. Porro explicabo temporibus, sunt officia vel corporis."},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Quae dolores consequatur officia [[at]] rem voluptatibus hic eos quo commodi minima? Ipsum commodi quae dolorum eum repudiandae provident saepe reiciendis vitae."},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Voluptate porro obcaecati numquam [[asperiores]] dignissimos accusantium veniam, dolore, illo doloremque dolorum, qui aliquam repudiandae nihil neque autem. Quaerat mollitia error molestias."},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Minima natus temporibus [[minus]] officia, repudiandae asperiores, nemo cum fugit rem repellat sequi iusto vero, explicabo corporis! Perspiciatis sequi, consectetur eveniet voluptate!"},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Optio rerum nam architecto [[mollitia]] unde officia tempora reiciendis omnis quas, expedita, quia ad culpa. Porro explicabo temporibus, sunt officia vel corporis."},
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Quae dolores consequatur officia [[at]] rem voluptatibus hic eos quo commodi minima? Ipsum commodi quae dolorum eum repudiandae provident saepe reiciendis vitae."},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Voluptate porro obcaecati numquam [[asperiores]] dignissimos accusantium veniam, dolore, illo doloremque dolorum, qui aliquam repudiandae nihil neque autem. Quaerat mollitia error molestias."},
		&Paragraph{"Lorem [[ipsum]] dolor sit amet, consectetur adipisicing elit. Minima natus temporibus [[minus]] officia, repudiandae asperiores, nemo cum fugit rem repellat sequi iusto vero, explicabo corporis! Perspiciatis sequi, consectetur eveniet voluptate!"},
	},
}

var Second = Page{
	URL:   "http://localhost:8000/second",
	Title: "Second",
	Story: []Item{
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Minus dolor ab aliquam sunt sit, eaque animi ut recusandae! Vitae ad rem eaque accusantium ex distinctio temporibus quo tempore? Vero, deserunt."},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. In consectetur nesciunt natus eos quos voluptates debitis consequuntur! Maxime quod libero ipsa sed, nihil at, quam rem consectetur corrupti sequi corporis."},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Neque doloremque eum amet atque vel consectetur veritatis dolore necessitatibus nostrum voluptatem facilis animi nulla dolorem tempore illum molestias cupiditate, pariatur, enim."},
		&Paragraph{"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Omnis placeat facilis enim aliquam aliquid voluptatum sequi accusantium repellat nihil nemo quasi excepturi mollitia explicabo quas, minus, corporis voluptatem pariatur in?"},
	},
}
