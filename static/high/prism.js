/* PrismJS 1.10.0
http://prismjs.com/download.html?themes=prism&languages=markup+css+clike+javascript+jass&plugins=line-numbers */
var _self = (typeof window !== 'undefined')
	? window   // if in browser
	: (
		(typeof WorkerGlobalScope !== 'undefined' && self instanceof WorkerGlobalScope)
		? self // if in worker
		: {}   // if in node js
	);

/**
 * Prism: Lightweight, robust, elegant syntax highlighting
 * MIT license http://www.opensource.org/licenses/mit-license.php/
 * @author Lea Verou http://lea.verou.me
 */

var Prism = (function(){

// Private helper vars
var lang = /\blang(?:uage)?-(\w+)\b/i;
var uniqueId = 0;

var _ = _self.Prism = {
	manual: _self.Prism && _self.Prism.manual,
	disableWorkerMessageHandler: _self.Prism && _self.Prism.disableWorkerMessageHandler,
	util: {
		encode: function (tokens) {
			if (tokens instanceof Token) {
				return new Token(tokens.type, _.util.encode(tokens.content), tokens.alias);
			} else if (_.util.type(tokens) === 'Array') {
				return tokens.map(_.util.encode);
			} else {
				return tokens.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/\u00a0/g, ' ');
			}
		},

		type: function (o) {
			return Object.prototype.toString.call(o).match(/\[object (\w+)\]/)[1];
		},

		objId: function (obj) {
			if (!obj['__id']) {
				Object.defineProperty(obj, '__id', { value: ++uniqueId });
			}
			return obj['__id'];
		},

		// Deep clone a language definition (e.g. to extend it)
		clone: function (o) {
			var type = _.util.type(o);

			switch (type) {
				case 'Object':
					var clone = {};

					for (var key in o) {
						if (o.hasOwnProperty(key)) {
							clone[key] = _.util.clone(o[key]);
						}
					}

					return clone;

				case 'Array':
					return o.map(function(v) { return _.util.clone(v); });
			}

			return o;
		}
	},

	languages: {
		extend: function (id, redef) {
			var lang = _.util.clone(_.languages[id]);

			for (var key in redef) {
				lang[key] = redef[key];
			}

			return lang;
		},

		/**
		 * Insert a token before another token in a language literal
		 * As this needs to recreate the object (we cannot actually insert before keys in object literals),
		 * we cannot just provide an object, we need anobject and a key.
		 * @param inside The key (or language id) of the parent
		 * @param before The key to insert before. If not provided, the function appends instead.
		 * @param insert Object with the key/value pairs to insert
		 * @param root The object that contains `inside`. If equal to Prism.languages, it can be omitted.
		 */
		insertBefore: function (inside, before, insert, root) {
			root = root || _.languages;
			var grammar = root[inside];

			if (arguments.length == 2) {
				insert = arguments[1];

				for (var newToken in insert) {
					if (insert.hasOwnProperty(newToken)) {
						grammar[newToken] = insert[newToken];
					}
				}

				return grammar;
			}

			var ret = {};

			for (var token in grammar) {

				if (grammar.hasOwnProperty(token)) {

					if (token == before) {

						for (var newToken in insert) {

							if (insert.hasOwnProperty(newToken)) {
								ret[newToken] = insert[newToken];
							}
						}
					}

					ret[token] = grammar[token];
				}
			}

			// Update references in other language definitions
			_.languages.DFS(_.languages, function(key, value) {
				if (value === root[inside] && key != inside) {
					this[key] = ret;
				}
			});

			return root[inside] = ret;
		},

		// Traverse a language definition with Depth First Search
		DFS: function(o, callback, type, visited) {
			visited = visited || {};
			for (var i in o) {
				if (o.hasOwnProperty(i)) {
					callback.call(o, i, o[i], type || i);

					if (_.util.type(o[i]) === 'Object' && !visited[_.util.objId(o[i])]) {
						visited[_.util.objId(o[i])] = true;
						_.languages.DFS(o[i], callback, null, visited);
					}
					else if (_.util.type(o[i]) === 'Array' && !visited[_.util.objId(o[i])]) {
						visited[_.util.objId(o[i])] = true;
						_.languages.DFS(o[i], callback, i, visited);
					}
				}
			}
		}
	},
	plugins: {},

	highlightAll: function(async, callback) {
		_.highlightAllUnder(document, async, callback);
	},

	highlightAllUnder: function(container, async, callback) {
		var env = {
			callback: callback,
			selector: 'code[class*="language-"], [class*="language-"] code, code[class*="lang-"], [class*="lang-"] code'
		};

		_.hooks.run("before-highlightall", env);

		var elements = env.elements || container.querySelectorAll(env.selector);

		for (var i=0, element; element = elements[i++];) {
			_.highlightElement(element, async === true, env.callback);
		}
	},

	highlightElement: function(element, async, callback) {
		// Find language
		var language, grammar, parent = element;

		while (parent && !lang.test(parent.className)) {
			parent = parent.parentNode;
		}

		if (parent) {
			language = (parent.className.match(lang) || [,''])[1].toLowerCase();
			grammar = _.languages[language];
		}

		// Set language on the element, if not present
		element.className = element.className.replace(lang, '').replace(/\s+/g, ' ') + ' language-' + language;

		if (element.parentNode) {
			// Set language on the parent, for styling
			parent = element.parentNode;

			if (/pre/i.test(parent.nodeName)) {
				parent.className = parent.className.replace(lang, '').replace(/\s+/g, ' ') + ' language-' + language;
			}
		}

		var code = element.textContent;

		var env = {
			element: element,
			language: language,
			grammar: grammar,
			code: code
		};

		_.hooks.run('before-sanity-check', env);

		if (!env.code || !env.grammar) {
			if (env.code) {
				_.hooks.run('before-highlight', env);
				env.element.textContent = env.code;
				_.hooks.run('after-highlight', env);
			}
			_.hooks.run('complete', env);
			return;
		}

		_.hooks.run('before-highlight', env);

		if (async && _self.Worker) {
			var worker = new Worker(_.filename);

			worker.onmessage = function(evt) {
				env.highlightedCode = evt.data;

				_.hooks.run('before-insert', env);

				env.element.innerHTML = env.highlightedCode;

				callback && callback.call(env.element);
				_.hooks.run('after-highlight', env);
				_.hooks.run('complete', env);
			};

			worker.postMessage(JSON.stringify({
				language: env.language,
				code: env.code,
				immediateClose: true
			}));
		}
		else {
			env.highlightedCode = _.highlight(env.code, env.grammar, env.language);

			_.hooks.run('before-insert', env);

			env.element.innerHTML = env.highlightedCode;

			callback && callback.call(element);

			_.hooks.run('after-highlight', env);
			_.hooks.run('complete', env);
		}
	},

	highlight: function (text, grammar, language) {
		var tokens = _.tokenize(text, grammar);
		return Token.stringify(_.util.encode(tokens), language);
	},

	matchGrammar: function (text, strarr, grammar, index, startPos, oneshot, target) {
		var Token = _.Token;

		for (var token in grammar) {
			if(!grammar.hasOwnProperty(token) || !grammar[token]) {
				continue;
			}

			if (token == target) {
				return;
			}

			var patterns = grammar[token];
			patterns = (_.util.type(patterns) === "Array") ? patterns : [patterns];

			for (var j = 0; j < patterns.length; ++j) {
				var pattern = patterns[j],
					inside = pattern.inside,
					lookbehind = !!pattern.lookbehind,
					greedy = !!pattern.greedy,
					lookbehindLength = 0,
					alias = pattern.alias;

				if (greedy && !pattern.pattern.global) {
					// Without the global flag, lastIndex won't work
					var flags = pattern.pattern.toString().match(/[imuy]*$/)[0];
					pattern.pattern = RegExp(pattern.pattern.source, flags + "g");
				}

				pattern = pattern.pattern || pattern;

				// Don’t cache length as it changes during the loop
				for (var i = index, pos = startPos; i < strarr.length; pos += strarr[i].length, ++i) {

					var str = strarr[i];

					if (strarr.length > text.length) {
						// Something went terribly wrong, ABORT, ABORT!
						return;
					}

					if (str instanceof Token) {
						continue;
					}

					pattern.lastIndex = 0;

					var match = pattern.exec(str),
					    delNum = 1;

					// Greedy patterns can override/remove up to two previously matched tokens
					if (!match && greedy && i != strarr.length - 1) {
						pattern.lastIndex = pos;
						match = pattern.exec(text);
						if (!match) {
							break;
						}

						var from = match.index + (lookbehind ? match[1].length : 0),
						    to = match.index + match[0].length,
						    k = i,
						    p = pos;

						for (var len = strarr.length; k < len && (p < to || (!strarr[k].type && !strarr[k - 1].greedy)); ++k) {
							p += strarr[k].length;
							// Move the index i to the element in strarr that is closest to from
							if (from >= p) {
								++i;
								pos = p;
							}
						}

						/*
						 * If strarr[i] is a Token, then the match starts inside another Token, which is invalid
						 * If strarr[k - 1] is greedy we are in conflict with another greedy pattern
						 */
						if (strarr[i] instanceof Token || strarr[k - 1].greedy) {
							continue;
						}

						// Number of tokens to delete and replace with the new match
						delNum = k - i;
						str = text.slice(pos, p);
						match.index -= pos;
					}

					if (!match) {
						if (oneshot) {
							break;
						}

						continue;
					}

					if(lookbehind) {
						lookbehindLength = match[1].length;
					}

					var from = match.index + lookbehindLength,
					    match = match[0].slice(lookbehindLength),
					    to = from + match.length,
					    before = str.slice(0, from),
					    after = str.slice(to);

					var args = [i, delNum];

					if (before) {
						++i;
						pos += before.length;
						args.push(before);
					}

					var wrapped = new Token(token, inside? _.tokenize(match, inside) : match, alias, match, greedy);

					args.push(wrapped);

					if (after) {
						args.push(after);
					}

					Array.prototype.splice.apply(strarr, args);

					if (delNum != 1)
						_.matchGrammar(text, strarr, grammar, i, pos, true, token);

					if (oneshot)
						break;
				}
			}
		}
	},

	tokenize: function(text, grammar, language) {
		var strarr = [text];

		var rest = grammar.rest;

		if (rest) {
			for (var token in rest) {
				grammar[token] = rest[token];
			}

			delete grammar.rest;
		}

		_.matchGrammar(text, strarr, grammar, 0, 0, false);

		return strarr;
	},

	hooks: {
		all: {},

		add: function (name, callback) {
			var hooks = _.hooks.all;

			hooks[name] = hooks[name] || [];

			hooks[name].push(callback);
		},

		run: function (name, env) {
			var callbacks = _.hooks.all[name];

			if (!callbacks || !callbacks.length) {
				return;
			}

			for (var i=0, callback; callback = callbacks[i++];) {
				callback(env);
			}
		}
	}
};

var Token = _.Token = function(type, content, alias, matchedStr, greedy) {
	this.type = type;
	this.content = content;
	this.alias = alias;
	// Copy of the full string this token was created from
	this.length = (matchedStr || "").length|0;
	this.greedy = !!greedy;
};

Token.stringify = function(o, language, parent) {
	if (typeof o == 'string') {
		return o;
	}

	if (_.util.type(o) === 'Array') {
		return o.map(function(element) {
			return Token.stringify(element, language, o);
		}).join('');
	}

	var env = {
		type: o.type,
		content: Token.stringify(o.content, language, parent),
		tag: 'span',
		classes: ['token', o.type],
		attributes: {},
		language: language,
		parent: parent
	};

	if (o.alias) {
		var aliases = _.util.type(o.alias) === 'Array' ? o.alias : [o.alias];
		Array.prototype.push.apply(env.classes, aliases);
	}

	_.hooks.run('wrap', env);

	var attributes = Object.keys(env.attributes).map(function(name) {
		return name + '="' + (env.attributes[name] || '').replace(/"/g, '&quot;') + '"';
	}).join(' ');

	return '<' + env.tag + ' class="' + env.classes.join(' ') + '"' + (attributes ? ' ' + attributes : '') + '>' + env.content + '</' + env.tag + '>';

};

if (!_self.document) {
	if (!_self.addEventListener) {
		// in Node.js
		return _self.Prism;
	}

	if (!_.disableWorkerMessageHandler) {
		// In worker
		_self.addEventListener('message', function (evt) {
			var message = JSON.parse(evt.data),
				lang = message.language,
				code = message.code,
				immediateClose = message.immediateClose;

			_self.postMessage(_.highlight(code, _.languages[lang], lang));
			if (immediateClose) {
				_self.close();
			}
		}, false);
	}

	return _self.Prism;
}

//Get current script and highlight
var script = document.currentScript || [].slice.call(document.getElementsByTagName("script")).pop();

if (script) {
	_.filename = script.src;

	if (!_.manual && !script.hasAttribute('data-manual')) {
		if(document.readyState !== "loading") {
			if (window.requestAnimationFrame) {
				window.requestAnimationFrame(_.highlightAll);
			} else {
				window.setTimeout(_.highlightAll, 16);
			}
		}
		else {
			document.addEventListener('DOMContentLoaded', _.highlightAll);
		}
	}
}

return _self.Prism;

})();

if (typeof module !== 'undefined' && module.exports) {
	module.exports = Prism;
}

// hack for components to work correctly in node.js
if (typeof global !== 'undefined') {
	global.Prism = Prism;
}
;
Prism.languages.jass = {
	'comment': [
		{
			pattern: /(^|[^\\:])\/\/.*/,
			lookbehind: true
		},
		{
			pattern: /(^|[^\\])\/\*[\s\S]*?(?:\*\/|$)/,
			lookbehind: true
		}
	],
	'value': [
		 /\b(?:null|true|false)\b/,
		 /(?:\d*\.|\$|0x|\b)[0-9a-f]+(?!')\b/,
		 /'[A-Za-z0-9]{4}'/i,
		{
			pattern: /(["])(?:\\(?:\r\n|[\s\S])|(?!\1)[^\\\r\n])*\1/,
			greedy: true
		},
		/\s(?=function)\s\w+/
	],
	'type': /\b(?:nothing|array|boolean|integer|real|string|trigger|timer|unit|item)(?:\$|\b)/,
	'native': /\b(?:A(?:cos|bilityId(?:|2String)|d(?:justCameraField|d(?:I(?:ndicator|temTo(?:Stock|AllStock))|Lightning(?:Ex|)|ResourceAmount|Spe(?:cialEffect(?:Loc||Target)|llEffect(?:Loc||Target(?:|ById)|ById(?:Loc|)))|HeroXP|Unit(?:AnimationProperties|To(?:Stock|AllStock))|WeatherEffect|PlayerTechResearched))|t(?:an(?:2|)|tachSoundToUnit)|nd|sin)|Kill(?:Destructable|SoundWhenDone|Unit)|E(?:xecuteFunc|n(?:d(?:Game|ThematicMusic|CinematicScene)|able(?:Occlusion|PreSelect|Trigger|User(?:UI|Control)|Select|DragSelect|W(?:orldFogBoundary|eatherEffect)|MinimapFilterButtons)|um(?:DestructablesInRect|ItemsInRect)))|D(?:i(?:s(?:play(?:T(?:extToPlayer|imedText(?:FromPlayer|ToPlayer))|LoadDialog|CineFilter)|able(?:Trigger|RestartMission))|alog(?:D(?:estroy|isplay)|SetMessage|C(?:lear|reate)|Add(?:QuitButton|Button)))|oNotSaveReplay|e(?:cUnitAbilityLevel|g2Rad|str(?:oy(?:DefeatCondition|U(?:nitPool|bersplat)|T(?:extTag|rigger|imer(?:Dialog|))|F(?:o(?:gModifier|rce)|ilter)|I(?:temPool|mage)|Effect|Group|BoolExpr|Quest|L(?:ightning|eaderboard)|Multiboard|Condition)|uctableRestoreLife)|f(?:ineStartLocation(?:Loc|)|eatConditionSetDescription)))|G(?:et(?:A(?:bility(?:Sound(?:|ById)|Effect(?:|ById))|llyColorFilterState|ttacker|IDifficulty)|C(?:ustomCampaignButtonVisible|licked(?:Button|Dialog)|hangingUnit(?:PrevOwner|)|re(?:atureDensity|epCampFilterState)|a(?:ncelledStructure|mera(?:BoundM(?:in(?:X|Y)|ax(?:X|Y))|TargetPosition(?:Z|Loc|X|Y)|Field|EyePosition(?:Z|Loc|X|Y)|Margin))|onstruct(?:ingStructure|edStructure))|E(?:vent(?:Player(?:State|ChatString(?:Matched|))|GameState|TargetUnit|D(?:amage(?:Source|)|etectingPlayer)|UnitState)|xpiredTimer|n(?:um(?:Destructable|Item|Player|Unit)|teringUnit))|Manipulat(?:edItem|ingUnit)|BuyingUnit|F(?:ilter(?:Destructable|Item|Player|Unit)|loatGameState|ood(?:Used|Made))|I(?:ntegerGameState|ssuedOrderId|tem(?:Y|Type(?:Id|)|X|Name|Level|UserData|Player|Charges))|H(?:ero(?:ProperName|S(?:killPoints|tr)|Level|XP|Agi|Int)|andleId)|KillingUnit|O(?:wningPlayer|bjectName|rder(?:Target(?:Destructable||Item|Unit)|edUnit|Point(?:Loc|X|Y)))|S(?:pell(?:Ability(?:|Id|Unit)|Target(?:Destructable|Loc|Item|Unit))|ummon(?:ingUnit|edUnit)|aveBasicFilename|t(?:artLoc(?:ation(?:Loc|X|Y)|Prio(?:Slot|))|ored(?:String|Boolean|Integer|Real))|o(?:und(?:Is(?:Loading|Playing)|Duration|FileDuration)|ld(?:Item|Unit))|ellingUnit)|R(?:e(?:ct(?:Center(?:X|Y)|M(?:in(?:X|Y)|ax(?:X|Y)))|s(?:ource(?:Density|Amount)|earch(?:ingUnit|ed)|cuer)|viv(?:ingUnit|ableUnit))|andom(?:Int|Real))|T(?:e(?:ams|rrain(?:Type|CliffLevel|Variance))|ournament(?:Score|Finish(?:Now(?:Player|Rule)|SoonTimeRemaining))|imeOfDayScale|r(?:a(?:nsportUnit|inedUnit(?:Type|))|igger(?:Destructable|ing(?:Tr(?:igger|ackable)|Region)|E(?:xecCount|v(?:entId|alCount))|Player|Widget|Unit)))|W(?:i(?:nningPlayer|dget(?:Life|X|Y))|orldBounds)|Player(?:Race|S(?:t(?:a(?:rtLocation|te)|ructureCount)|electable|lotState|core)|Name|T(?:ypedUnitCount|axRate|e(?:am|ch(?:Count|MaxAllowed|Researched)))|s|UnitCount|Alliance|Co(?:ntroller|lor)|Handicap(?:|XP)|Id)|L(?:e(?:a(?:vingUnit|rn(?:edSkill(?:Level|)|ingUnit))|velingUnit)|o(?:adedUnit|ca(?:l(?:Player|ized(?:String|Hotkey))|tion(?:Z|X|Y)))|ightningColor(?:R|G|B|A))|Game(?:Speed|Difficulty|Placement|TypeSelected)|D(?:e(?:cayingUnit|tectedUnit|structable(?:TypeId|X|Name|MaxLife|Life|OccluderHeight|Y)|faultDifficulty)|yingUnit)|Unit(?:A(?:cquireRange|bilityLevel)|P(?:ropWindow|ointValue(?:ByType|))|MoveSpeed|Name|CurrentOrder|UserData|F(?:acing|lyHeight|ood(?:Used|Made))|Ra(?:ce|lly(?:Destructable|Point|Unit))|Default(?:FlyHeight|MoveSpeed|AcquireRange|TurnSpeed|PropWindow)|Y|State|L(?:evel|oc)|X|T(?:urnSpeed|ypeId)))|roup(?:ImmediateOrder(?:|ById)|RemoveUnit|PointOrder(?:Loc||ById(?:Loc|))|Clear|EnumUnits(?:Of(?:Player|Type(?:Counted|))|Selected|InR(?:ange(?:OfLoc(?:Counted|)||Counted)|ect(?:Counted|)))|AddUnit|TargetOrder(?:|ById)))|F(?:i(?:nishUbersplat|rstOfGroup|lter)|l(?:ush(?:Stored(?:Integer|Real|Mission|String|Boolean|Unit)|ParentHashtable|ChildHashtable|GameCache)|ashQuestDialogButton)|o(?:g(?:M(?:askEnable|odifierSt(?:op|art))|Enable)|r(?:Force|Group|ce(?:PlayerStartLocation|C(?:inematicSubtitles|lear|ampaignSelectScreen)|QuestDialogUpdate|RemovePlayer|Enum(?:Players(?:Counted|)|Allies|Enemies)|UI(?:Cancel|Key)|AddPlayer))))|I(?:n(?:it(?:GameCache|Hashtable)|cUnitAbilityLevel)|temPool(?:AddItemType|RemoveItemType)|s(?:SuspendedXP|P(?:oint(?:InRegion|Blighted)|layer(?:RacePrefSet|InForce|Observer|Ally|Enemy))|GameTypeSupported|sue(?:I(?:mmediateOrder(?:|ById)|nstant(?:TargetOrder(?:|ById)|PointOrder(?:|ById)))|PointOrder(?:Loc||ById(?:Loc|))|TargetOrder(?:|ById)|Neutral(?:TargetOrder(?:|ById)|ImmediateOrder(?:|ById)|PointOrder(?:|ById))|BuildOrder(?:|ById))|No(?:VictoryCheat|DefeatCheat)|T(?:errainPathable|imerDialogDisplayed|rigger(?:Enabled|WaitOnSleeps))|Item(?:I(?:nvulnerable|d(?:Sellable|P(?:owerup|awnable)))|P(?:owerup|awnable)|Visible|Owned|Sellable)|Unit(?:Masked|Detected|Selected||OwnedByPlayer|Paused|Enemy|I(?:llusion|dType|n(?:Group|Transport|visible|R(?:ange(?:Loc||XY)|egion)|Force))|Hidden|Visible|Type|Race|Ally|Loaded|Fogged)|HeroUnitId|Fog(?:MaskEnabled|gedToPlayer|Enabled)|VisibleToPlayer|CineFilterDisplayed|M(?:a(?:pFlagSet|skedToPlayer)|ultiboard(?:Minimized|Displayed))|L(?:eaderboardDisplayed|ocation(?:FoggedToPlayer|InRegion|VisibleToPlayer|MaskedToPlayer))|Quest(?:Required|Completed|ItemCompleted|Failed|Enabled|Discovered)|DestructableInvulnerable)|2(?:R|S))|Or(?:|derId(?:|2String))|M(?:ove(?:RectTo(?:Loc|)|L(?:ightning(?:Ex|)|ocation))|ultiboard(?:Minimize|S(?:et(?:ColumnCount|TitleText(?:Color|)|Item(?:s(?:Value(?:Color|)|Style|Icon|Width)|Value(?:Color|)|Style|Icon|Width)|RowCount)|uppressDisplay)|Clear|Display|ReleaseItem|Get(?:ColumnCount|TitleText|Item|RowCount)))|L(?:o(?:ad(?:D(?:ialogHandle|e(?:structableHandle|featConditionHandle))|PlayerHandle|S(?:oundHandle|tr)|B(?:oolean(?:ExprHandle|)|uttonHandle)|U(?:bersplatHandle|nit(?:PoolHandle|Handle))|T(?:extTagHandle|imer(?:DialogHandle|Handle)|r(?:ackableHandle|igger(?:Handle|EventHandle|ActionHandle|ConditionHandle)))|G(?:roupHandle|ame)|AbilityHandle|I(?:nteger|mageHandle|tem(?:Handle|PoolHandle))|Re(?:al|ctHandle|gionHandle)|Multiboard(?:ItemHandle|Handle)|EffectHandle|Fo(?:g(?:StateHandle|ModifierHandle)|rceHandle)|L(?:ightningHandle|eaderboardHandle|ocationHandle)|Quest(?:ItemHandle|Handle)|WidgetHandle)|cation)|eaderboard(?:Remove(?:Item|PlayerItem)|S(?:ortItemsBy(?:Label|Value|Player)|et(?:ValueColor|Item(?:Value(?:Color|)|Style|Label(?:Color|))|S(?:tyle|izeByItemCount)|Label(?:Color|)))|Clear|Display|AddItem|HasPlayerItem|Get(?:LabelText|PlayerIndex|ItemCount)))|N(?:ewSoundEnvironment|ot)|HaveS(?:aved(?:Integer|Real|Handle|Boolean|String)|tored(?:Integer|Real|String|Boolean|Unit))|P(?:a(?:nCameraTo(?:Timed(?:WithZ|)|WithZ|)|use(?:Game|Timer|CompAI|Unit))|ow|ingMinimap(?:Ex|)|la(?:ceRandom(?:Item|Unit)|y(?:M(?:odelCinematic|usic(?:Ex|))|er(?:|SetLeaderboard|GetLeaderboard)|ThematicMusic(?:Ex|)|Cinematic))|reload(?:|er|End(?:Ex|)|Refresh|Gen(?:Clear|End|Start)|Start))|S(?:2(?:R|I)|u(?:spend(?:TimeOfDay|HeroXP)|bString)|t(?:art(?:CampaignAI|MeleeAI|Sound)|o(?:re(?:Integer|Real|String|Boolean|Unit)|p(?:Sound|Camera|Music))|ring(?:Hash|Case|Length))|quareRoot|yncS(?:tored(?:Integer|Real|String|Boolean|Unit)|elections)|ave(?:D(?:ialogHandle|e(?:structableHandle|featConditionHandle))|PlayerHandle|S(?:oundHandle|tr)|B(?:oolean(?:ExprHandle|)|uttonHandle)|U(?:bersplatHandle|nit(?:PoolHandle|Handle))|T(?:extTagHandle|imer(?:DialogHandle|Handle)|r(?:ackableHandle|igger(?:Handle|EventHandle|ActionHandle|ConditionHandle)))|G(?:roupHandle|ame(?:Exists||Cache))|L(?:ightningHandle|eaderboardHandle|ocationHandle)|I(?:nteger|mageHandle|tem(?:Handle|PoolHandle))|Fo(?:g(?:StateHandle|ModifierHandle)|rceHandle)|Multiboard(?:ItemHandle|Handle)|Re(?:al|ctHandle|gionHandle)|EffectHandle|AbilityHandle|Quest(?:ItemHandle|Handle)|WidgetHandle)|in|how(?:U(?:bersplat|nit)|I(?:mage|nterface)|Destructable)|e(?:lect(?:Unit|HeroSkill)|t(?:Hero(?:Int|Level|XP|Agi|Str)|C(?:re(?:atureDensity|epCampFilterState)|ustomCampaignButtonVisible|ine(?:matic(?:Camera|Scene)|Filter(?:BlendMode|End(?:Color|UV)|Start(?:Color|UV)|Duration|Tex(?:MapFlags|ture)))|am(?:era(?:Position|Bounds|TargetController|OrientController|QuickPosition|RotateMode|Field)|paign(?:MenuRace(?:Ex|)|Available)))|D(?:e(?:structable(?:Invulnerable|Animation(?:Speed|)|MaxLife|Life|OccluderHeight)|faultDifficulty)|oodadAnimation(?:Rect|)|ayNightModels)|F(?:loatGameState|ogStateR(?:adius(?:Loc|)|ect))|I(?:mage(?:Type|Position|Co(?:nstantHeight|lor)|AboveWater|Render(?:Always|))|nt(?:egerGameState|roShot(?:Model|Text))|tem(?:P(?:awnable|layer|osition)|Visible|Drop(?:ID|pable|OnDeath)|Charges|UserData|Invulnerable|TypeSlots))|M(?:usic(?:PlayPosition|Volume)|issionAvailable|ap(?:Flag|Description|Music|Name))|Blight(?:Loc||Point|Rect)|S(?:ound(?:D(?:uration|istance(?:s|Cutoff))|P(?:aramsFromLabel|layPosition|itch|osition)|C(?:one(?:Angles|Orientation)|hannel)|V(?:elocity|olume))|tartLocPrio(?:Count|)|kyModel)|R(?:e(?:ct(?:FromLoc|)|s(?:ource(?:Density|Amount)|ervedLocalHeroButtons))|andomSeed)|U(?:bersplatRender(?:Always|)|nit(?:A(?:cquireRange|nimation(?:WithRarity||ByIndex)|bilityLevel)|Invulnerable|C(?:olor|reepGuard)|BlendTime|F(?:acing(?:|Timed)|lyHeight|og)|Rescu(?:able|eRange)|S(?:cale|tate)|T(?:urnSpeed|ypeSlots|imeScale)|MoveSpeed|P(?:ropWindow|athing|osition(?:Loc|))|Owner|X|Y|LookAt|VertexColor|Exploded|Use(?:Food|rData)))|T(?:e(?:ams|xtTag(?:Color|Suspended|Fadepoint|Text|Lifespan|Age|P(?:ermanent|os(?:Unit|))|V(?:elocity|isibility))|rrain(?:Type|Pathable|Fog(?:Ex|)))|imeOfDayScale|utorialCleared|hematicMusicPlayPosition)|W(?:idgetLife|ater(?:BaseColor|Deforms))|EdCinematicAvailable|Player(?:A(?:bilityAvailable|lliance)|s|Name|T(?:e(?:am|ch(?:MaxAllowed|Researched))|axRate)|Race(?:Selectable|Preference)|Sta(?:rtLocation|te)|UnitsOwner|OnScoreScreen|Handicap(?:|XP)|Co(?:ntroller|lor))|Al(?:tMinimapIcon|l(?:yColorFilterState|ItemTypeSlots|UnitTypeSlots))|Game(?:Difficulty|Speed|Placement|TypeSupported)|OpCinematicAvailable|LightningColor)))|R(?:e(?:viveHero(?:Loc|)|nameSaveDirectory|c(?:ycleGuardPosition|t(?:FromLoc|))|s(?:et(?:U(?:nitLookAt|bersplat)|T(?:rigger|errainFog|oGameCamera))|t(?:artGame|oreUnit)|ume(?:Music|Timer))|move(?:GuardPosition|WeatherEffect|Location|Save(?:Directory|d(?:Integer|Real|Handle|Boolean|String))|AllGuardPositions|Unit(?:|From(?:Stock|AllStock))|Destructable|Item(?:|From(?:Stock|AllStock))|Re(?:gion|ct)|Player)|gi(?:on(?:Clear(?:Rect|Cell(?:AtLoc|))|Add(?:Rect|Cell(?:AtLoc|)))|sterStackedSound)|loadGame(?:CachesFromDisk|))|ad2Deg|2(?:S(?:W|)|I))|T(?:an|errainDeform(?:Crater|Wave|Stop(?:All|)|R(?:ipple|andom))|imer(?:Dialog(?:Display|Set(?:Ti(?:meColor|tle(?:Color|))|RealTimeRemaining|Speed))|Get(?:Elapsed|Remaining|Timeout)|Start)|rigger(?:Clear(?:Actions|Conditions)|S(?:ync(?:Start|Ready)|leepAction)|Re(?:move(?:Action|Condition)|gister(?:T(?:imerE(?:xpireEvent|vent)|rackable(?:TrackEvent|HitEvent))|Unit(?:InRange|StateEvent|Event)|EnterRegion|Player(?:ChatEvent|Event|StateEvent|UnitEvent|AllianceChange)|D(?:ialog(?:Event|ButtonEvent)|eathEvent)|LeaveRegion|Game(?:Event|StateEvent)|FilterUnitEvent|VariableEvent))|Add(?:Action|Condition)|E(?:valuate|xecute(?:Wait|))|Wait(?:OnSleeps|ForSound)))|V(?:ersion(?:Compatible|Supported|Get)|olumeGroup(?:SetVolume|Reset))|Que(?:st(?:CreateItem|ItemSet(?:Completed|Description)|Set(?:Required|Completed|IconPath|Failed|D(?:iscovered|escription)|Enabled|Title))|ue(?:UnitAnimation|DestructableAnimation))|C(?:a(?:meraSet(?:TargetNoise(?:Ex|)|S(?:moothingFactor|ourceNoise(?:Ex|))|up(?:Apply(?:|WithZ|ForceDuration(?:WithZ|))|Set(?:DestPosition|Field)|Get(?:Field|DestPosition(?:Loc|X|Y))))|chePlayerHeroData)|lear(?:MapMusic|Selection|TextMessages)|r(?:eate(?:C(?:ameraSetup|orpse)|U(?:bersplat|nit(?:Pool||AtLoc(?:ByName|)|ByName))|T(?:extTag|imer(?:Dialog|)|r(?:igger|ackable))|Group|Leaderboard|I(?:mage|tem(?:Pool|))|Fo(?:gModifierR(?:adius(?:Loc|)|ect)|rce)|De(?:adDestructable(?:Z|)|structable(?:Z|)|featCondition)|BlightedGoldmine|M(?:ultiboard|IDISound)|Quest|Sound(?:|F(?:ilenameWithLabel|romLabel))|Region)|ipplePlayer)|o(?:mmandAI|pySaveGame|s|n(?:dition|vert(?:A(?:llianceType|IDifficulty|ttackType)|TexMapFlags|D(?:ialogEvent|amageType)|Unit(?:State|Event|Type)|I(?:GameState|temType)|BlendMode|LimitOp|S(?:oundType|tartLocPrio)|CameraField|EffectType|W(?:eaponType|idgetEvent)|V(?:ersion|olumeGroup)|Ra(?:rityControl|ce(?:Pref|))|F(?:ogState|GameState)|Game(?:Type|Speed|Event|Difficulty)|Map(?:Flag|Setting|Control|Density|Visibility)|P(?:athingType|la(?:cement|yer(?:Color|GameResult|Event|UnitEvent|S(?:tate|core|lotState)))))))|h(?:eat|angeLevel|ooseRandom(?:NPBuilding|Creep|Item(?:Ex|))))|Waygate(?:SetDestination|Activate|GetDestination(?:X|Y)|IsActive)|Un(?:it(?:A(?:pplyTimedLife|dd(?:I(?:ndicator|tem(?:ToSlotById||ById))|Type|Sleep(?:Perm|)|Ability))|P(?:auseTimedLife|ool(?:RemoveUnitType|AddUnitType))|S(?:et(?:U(?:sesAltIcon|pgradeProgress)|ConstructionProgress)|tripHeroLevel|hareVision|uspendDecay)|D(?:ropItem(?:Point|Slot|Target)|amage(?:Target|Point))|I(?:temInSlot|sSleeping|d(?:2String|)|gnoreAlarm(?:Toggled|)|nventorySize)|UseItem(?:Point|Target|)|C(?:anSleep(?:Perm|)|ountBuffsEx)|Re(?:move(?:Type|Buffs(?:Ex|)|Item(?:FromSlot|)|Ability)|setCooldown)|Has(?:Item|BuffsEx)|M(?:akeAbilityPermanent|odifySkillPoints)|WakeUp)|registerStackedSound))(?=\()/,
	'function': [
		{
			pattern: /\w+(?=\(|\stakes\b)/i
		},
		{
			pattern: /(\b(?:function|initializer)\s)([\w]+)\b/,
			lookbehind: true
		}
	],
	'library': {
		pattern: /((?:\b(?:library|interface|initializer|requires|struct)|#(?:define|ifdef|ifndef|include|undef))\s+)([\w\<\>\/\\\."]+)/,
		lookbehind: true,
		inside: {
			punctuation: /[ .\\]/
		}
	},
	'keyword': /(?:\b(?:and|array|call|constant|else|elseif|endfunction|endglobals|endif|endloop|exitwhen|extends|function|globals|then|if|local|loop|native|not|or|return|returns|set|takes|type|library|endlibrary|initializer|requires|public|private|struct|endstruct|method|endmethod|static|interface|endinterface|operator|debug)|\s*#(?:define|endif|ifdef|ifndef|include|undef))(?:\$|\b)/,
	'operator': /<[=>]?|>=?|[+\-*\/^=&]/i,
	'punctuation': /[,()\[\]]/
};
(function () {

	if (typeof self === 'undefined' || !self.Prism || !self.document) {
		return;
	}

	/**
	 * Plugin name which is used as a class name for <pre> which is activating the plugin
	 * @type {String}
	 */
	var PLUGIN_NAME = 'line-numbers';
	
	/**
	 * Regular expression used for determining line breaks
	 * @type {RegExp}
	 */
	var NEW_LINE_EXP = /\n(?!$)/g;

	/**
	 * Resizes line numbers spans according to height of line of code
	 * @param {Element} element <pre> element
	 */
	var _resizeElement = function (element) {
		var codeStyles = getStyles(element);
		var whiteSpace = codeStyles['white-space'];

		if (whiteSpace === 'pre-wrap' || whiteSpace === 'pre-line') {
			var codeElement = element.querySelector('code');
			var lineNumbersWrapper = element.querySelector('.line-numbers-rows');
			var lineNumberSizer = element.querySelector('.line-numbers-sizer');
			var codeLines = codeElement.textContent.split(NEW_LINE_EXP);

			if (!lineNumberSizer) {
				lineNumberSizer = document.createElement('span');
				lineNumberSizer.className = 'line-numbers-sizer';

				codeElement.appendChild(lineNumberSizer);
			}

			lineNumberSizer.style.display = 'block';

			codeLines.forEach(function (line, lineNumber) {
				lineNumberSizer.textContent = line || '\n';
				var lineSize = lineNumberSizer.getBoundingClientRect().height;
				lineNumbersWrapper.children[lineNumber].style.height = lineSize + 'px';
			});

			lineNumberSizer.textContent = '';
			lineNumberSizer.style.display = 'none';
		}
	};

	/**
	 * Returns style declarations for the element
	 * @param {Element} element
	 */
	var getStyles = function (element) {
		if (!element) {
			return null;
		}

		return window.getComputedStyle ? getComputedStyle(element) : (element.currentStyle || null);
	};

	window.addEventListener('resize', function () {
		Array.prototype.forEach.call(document.querySelectorAll('pre.' + PLUGIN_NAME), _resizeElement);
	});

	Prism.hooks.add('complete', function (env) {
		if (!env.code) {
			return;
		}

		// works only for <code> wrapped inside <pre> (not inline)
		var pre = env.element.parentNode;
		var clsReg = /\s*\bline-numbers\b\s*/;
		if (
			!pre || !/pre/i.test(pre.nodeName) ||
			// Abort only if nor the <pre> nor the <code> have the class
			(!clsReg.test(pre.className) && !clsReg.test(env.element.className))
		) {
			return;
		}

		if (env.element.querySelector('.line-numbers-rows')) {
			// Abort if line numbers already exists
			return;
		}

		if (clsReg.test(env.element.className)) {
			// Remove the class 'line-numbers' from the <code>
			env.element.className = env.element.className.replace(clsReg, ' ');
		}
		if (!clsReg.test(pre.className)) {
			// Add the class 'line-numbers' to the <pre>
			pre.className += ' line-numbers';
		}

		var match = env.code.match(NEW_LINE_EXP);
		var linesNum = match ? match.length + 1 : 1;
		var lineNumbersWrapper;

		var lines = new Array(linesNum + 1);
		lines = lines.join('<span></span>');

		lineNumbersWrapper = document.createElement('span');
		lineNumbersWrapper.setAttribute('aria-hidden', 'true');
		lineNumbersWrapper.className = 'line-numbers-rows';
		lineNumbersWrapper.innerHTML = lines;

		if (pre.hasAttribute('data-start')) {
			pre.style.counterReset = 'linenumber ' + (parseInt(pre.getAttribute('data-start'), 10) - 1);
		}

		env.element.appendChild(lineNumbersWrapper);

		_resizeElement(pre);

		Prism.hooks.run('line-numbers', env);
	});

	Prism.hooks.add('line-numbers', function (env) {
		env.plugins = env.plugins || {};
		env.plugins.lineNumbers = true;
	});
	
	/**
	 * Global exports
	 */
	Prism.plugins.lineNumbers = {
		/**
		 * Get node for provided line number
		 * @param {Element} element pre element
		 * @param {Number} number line number
		 * @return {Element|undefined}
		 */
		getLine: function (element, number) {
			if (element.tagName !== 'PRE' || !element.classList.contains(PLUGIN_NAME)) {
				return;
			}

			var lineNumberRows = element.querySelector('.line-numbers-rows');
			var lineNumberStart = parseInt(element.getAttribute('data-start'), 10) || 1;
			var lineNumberEnd = lineNumberStart + (lineNumberRows.children.length - 1);

			if (number < lineNumberStart) {
				number = lineNumberStart;
			}
			if (number > lineNumberEnd) {
				number = lineNumberEnd;
			}

			var lineIndex = number - lineNumberStart;

			return lineNumberRows.children[lineIndex];
		}
	};

}());
