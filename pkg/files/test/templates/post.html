<?php
Tokens::add(['#TITLE#' => Page::get_value('title')]);

$scripts = ['vue.min.js','vee-validate.min.js','axios.min.js','common.js','comments.js'];
$css = ['comments.css','pure.min.css'];
if (Session::$is_moderator)
{
  $scripts[] = 'content-tools.min.js';
  $scripts[] = 'cloudinary.js';
  $scripts[] = 'editor.js';
  $css[] = 'content-tools.min.css'; 
}
HTML::register_scripts($scripts);
HTML::script_meta_tags(['css' => $css]);

new RecentPosts();

include 'header.php';
?>
<div class="sidebar">
<h3>Recent Posts</h3>

#RECENT#

<h3>Archives</h3>

<?php Archive::get_content(); ?>

<h3>Categories</h3>

<?php (new Categories())->print_list(); ?>
</div>

<div class="content">
	<h1>#TITLE#</h1>
	<p class="meta">Published: <?php echo date('jS F Y', strtotime(Page::get_value('published'))); ?></p>

	<div data-editable data-name="main-content">
	<?php Page::insert_region_content(); ?>
	</div>

	<div>
		<div id="comments">
			<a name="comments"></a>

			<comment :comment="{ id: 0 }" :root="root"></comment>

			<comment v-for="comment in root.comments" :comment="comment" :root="root"></comment>
		</div>
	</div>
</div>

<?php include 'footer.php'; ?>