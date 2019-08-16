
var vueNews = new Vue({
	el: '#news',
	delimiters: ['@{', '}'],
	data: {
		currentArticlePage: 1,
		scenes: [],
	},
	methods: {
		getNextArticles: function() {
	        axios
	            .get('/news/pages/' + (vueNews.currentArticlePage + 1) + '/')
	                .then(function (response) {
						response.data.scenes.forEach(scene => {
							vueNews.scenes.push({
								// 作成日
			                    createdAt: scene.CreatedAtString,
			                    // 初心者おすすめ度
			                    recommenddedNumber: vueNews.getRecommendedNumberString(scene.RecommendedNumber),
			                    // タイトル（日本語）
			                    titleJp: scene.Movie.TitleJp,
			                    // タイトル（英語）
			                    titleEn: scene.Movie.TitleEn,
								// シーンタイトル
			                    sceneTitle: scene.Title,
			                    // シーン説明
			                    description: scene.Description,
							});
						});
						vueNews.currentArticlePage++
	                });
		},
		// 初心者おすすめ度の星の文字列を取得する
		getRecommendedNumberString: function(recommenddedNumber) {
			const START_BLACK = "★";
			const START_WHITE = "☆";
			var numberStriung = START_BLACK.repeat(recommenddedNumber);
			return numberStriung + START_WHITE.repeat(5 - recommenddedNumber);
		},
		// シーン詳細画面へ遷移
		moveToScene: function(scene) {
			window.location.href = "/scenes/" + scene.id + "/";
		},
	},
})
