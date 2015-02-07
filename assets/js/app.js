var Title = React.createClass({
	render() {
		return (
			<div class="page-header">
				<h1>Plantumlor</h1>
			</div>
		);
	},
});

var Editor = React.createClass({
	propTypes: {
		onChange:    React.PropTypes.func.isRequired,
		defaultText: React.PropTypes.string.isRequired,
	},
	componentDidMount() {
		this.props.behave = new Behave({
			textarea:   this.getDOMNode(),
			replaceTab: true,
			softTabs:   true,
			tabSize:    4,
			autoOpen:   true,
			overwrite:  true,
			autoStrip:  true,
			autoIndent: true,
			fence:      false
		});
	},
	componentWillUnmount() {
		this.props.destroy()
	},
	_onChange() {
		this.props.onChange(this.refs.textArea.getDOMNode().value)
	},
	render() {
		return (
			<textarea onChange={this._onChange} ref="textArea">{this.props.defaultText}</textarea>
		);
	},
});

var Image = React.createClass({
	propTypes: {
		text: React.PropTypes.string.isRequired,
	},
	fixedEncodeURIComponent (str) {
		return encodeURIComponent(str).replace(/[!'()]/g, escape).replace(/\*/g, "%2A");
	},
	render() {
		var url = "/transfer/" + Base64.btoa(RawDeflate.deflate(Base64.utob(this.props.text))).replace(/\//g, "_");
		return (
			<img src={url} />
		);
	},
});

var Content = React.createClass({
	propTypes: {
		defaultText: React.PropTypes.string.isRequired,
	},
	getInitialState() {
		return {
			newText: this.props.defaultText,
			oldText: "",
		};
	},
	change(text) {
		if (text !== "") {
			this.setState({
				newText: text,
			});
		}
	},
	rerender() {
		if (this.state.newText !== this.state.oldText) {
			this.setState({
				oldText: this.state.newText,
			});
			console.log("rerender")
		}
	},
	componentDidMount: function() {
		this.interval = setInterval(this.rerender, 1000);
	},
	componentWillUnmount: function() {
		clearInterval(this.interval);
	},
	render() {
		return (
			<div>
				<Title />
				<Editor onChange={this.change} defaultText={this.props.defaultText} />
				<hr />
				<Image text={this.state.oldText} />
			</div>
		);
	},
});

React.render(
	<Content defaultText="Bob->Alice : hello" />,
	document.getElementById('content')
);

